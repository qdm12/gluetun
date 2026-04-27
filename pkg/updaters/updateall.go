package updaters

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/qdm12/dns/v2/pkg/doh"
	dnsprovider "github.com/qdm12/dns/v2/pkg/provider"
	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/common"
	publicipapi "github.com/qdm12/gluetun/internal/publicip/api"
	"github.com/qdm12/gluetun/internal/updater/resolver"
	"github.com/qdm12/gluetun/internal/updater/unzip"
	"github.com/qdm12/gluetun/pkg/updaters/airvpn"
	"github.com/qdm12/gluetun/pkg/updaters/cyberghost"
	"github.com/qdm12/gluetun/pkg/updaters/expressvpn"
	"github.com/qdm12/gluetun/pkg/updaters/fastestvpn"
	"github.com/qdm12/gluetun/pkg/updaters/giganews"
	"github.com/qdm12/gluetun/pkg/updaters/hidemyass"
	"github.com/qdm12/gluetun/pkg/updaters/ipvanish"
	"github.com/qdm12/gluetun/pkg/updaters/ivpn"
	"github.com/qdm12/gluetun/pkg/updaters/mullvad"
	"github.com/qdm12/gluetun/pkg/updaters/nordvpn"
	"github.com/qdm12/gluetun/pkg/updaters/perfectprivacy"
	"github.com/qdm12/gluetun/pkg/updaters/privado"
	"github.com/qdm12/gluetun/pkg/updaters/privateinternetaccess"
	"github.com/qdm12/gluetun/pkg/updaters/privatevpn"
	"github.com/qdm12/gluetun/pkg/updaters/protonvpn"
	"github.com/qdm12/gluetun/pkg/updaters/purevpn"
	"github.com/qdm12/gluetun/pkg/updaters/slickvpn"
	"github.com/qdm12/gluetun/pkg/updaters/surfshark"
	"github.com/qdm12/gluetun/pkg/updaters/torguard"
	"github.com/qdm12/gluetun/pkg/updaters/vpnsecure"
	"github.com/qdm12/gluetun/pkg/updaters/vpnunlimited"
	"github.com/qdm12/gluetun/pkg/updaters/vyprvpn"
	"github.com/qdm12/gluetun/pkg/updaters/windscribe"
	"github.com/qdm12/gosettings"
)

// UpdateAllSettings contains the configuration for the [UpdateAll] function.
type UpdateAllSettings struct {
	// OutputPath is the directory where the provider JSON files will be written.
	// It defaults to the current directory if left unset.
	OutputPath *string
	// ProtonEmail is the email for the ProtonVPN account, which is required
	// to update ProtonVPN servers.
	ProtonEmail *string
	// ProtonPassword is the password for the ProtonVPN account, which is required
	// to update ProtonVPN servers.
	ProtonPassword *string
	// IpinfoToken is the API token for the IPInfo public IP API.
	// If not provided, the IP fetcher will still work but may be
	// subject to stricter rate limits for ipinfo.io.
	IpinfoToken *string
	// MinServers is a map of provider name to minimum number of servers required for a successful update.
	// If a provider has fewer servers than the specified minimum, the update will be considered a failure
	// for that provider. If the provider name is not found in this map, it is assumed that there is no
	// minimum server requirement for that provider.
	MinServers map[string]uint
}

func (s *UpdateAllSettings) setDefaults() {
	s.OutputPath = gosettings.DefaultPointer(s.OutputPath, "")
	s.ProtonEmail = gosettings.DefaultPointer(s.ProtonEmail, "")
	s.ProtonPassword = gosettings.DefaultPointer(s.ProtonPassword, "")
	s.IpinfoToken = gosettings.DefaultPointer(s.IpinfoToken, "")
	if s.MinServers == nil {
		s.MinServers = make(map[string]uint)
	}
}

var (
	errProtonEmailRequired        = errors.New("proton email is required for updating ProtonVPN servers")
	errProtonPasswordRequired     = errors.New("proton password is required for updating ProtonVPN servers")
	errMinServersProviderNotFound = errors.New("provider name in MinServers not found in list of all providers")
)

func (s *UpdateAllSettings) validate() error {
	switch {
	case *s.ProtonEmail == "":
		return fmt.Errorf("%w", errProtonEmailRequired)
	case *s.ProtonPassword == "":
		return fmt.Errorf("%w", errProtonPasswordRequired)
	}

	allProviders := providers.All()
	for providerName := range s.MinServers {
		if !slices.Contains(allProviders, providerName) {
			return fmt.Errorf("%w: %s", errMinServersProviderNotFound, providerName)
		}
	}

	return nil
}

type Logger interface {
	Warn(message string)
}

var errUpdateAllFailed = errors.New("update failed for one or more providers")

// UpdateAll fetches server data for all providers in parallel and writes each
// provider's data as <provider-name>.json into the output path specified.
// Errors for individual providers are collected and returned as a combined error.
func UpdateAll(ctx context.Context, client *http.Client, logger Logger, settings UpdateAllSettings) error {
	settings.setDefaults()
	err := settings.validate()
	if err != nil {
		return fmt.Errorf("validating settings: %w", err)
	}

	const permission = 0o755
	err = os.MkdirAll(*settings.OutputPath, permission)
	if err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	dohDialer, err := doh.New(doh.Settings{
		UpstreamResolvers: []dnsprovider.Provider{
			dnsprovider.Cloudflare(),
			dnsprovider.Google(),
		},
	})
	if err != nil {
		return fmt.Errorf("creating updater DoH dialer: %w", err)
	}
	parallelResolver := resolver.NewParallelResolver(dohDialer)
	unzipper := unzip.New(client)
	ipFetcher, err := buildIPFetcher(client, logger, *settings.IpinfoToken)
	if err != nil {
		return fmt.Errorf("creating IP fetcher: %w", err)
	}
	fetchers := buildFetchers(client, parallelResolver, unzipper, ipFetcher,
		logger, *settings.ProtonEmail, *settings.ProtonPassword)

	results := make(chan error)

	allProviders := providers.All()
	for _, providerName := range allProviders {
		go func(providerName string) {
			fetcher := fetchers[providerName]
			minServers := settings.MinServers[providerName]
			err := fetchAndWrite(ctx, providerName, fetcher, minServers, *settings.OutputPath)
			if err != nil {
				err = fmt.Errorf("provider %s: %w", providerName, err)
			}
			results <- err
		}(providerName)
	}

	errs := make([]string, 0, len(allProviders))
	for range allProviders {
		err := <-results
		if err != nil {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("%w:\n%s", errUpdateAllFailed, strings.Join(errs, "\n"))
	}
	return nil
}

func buildIPFetcher(client *http.Client, logger Logger, ipinfoToken string) (
	fetcher *publicipapi.ResilientFetcher, err error,
) {
	nameTokenPairs := []publicipapi.NameToken{
		{Name: string(publicipapi.IPInfo), Token: ipinfoToken},
		{Name: string(publicipapi.IP2Location)},
		{Name: string(publicipapi.IfConfigCo)},
	}
	fetchers, err := publicipapi.New(nameTokenPairs, client)
	if err != nil {
		return nil, fmt.Errorf("creating public IP fetchers: %w", err)
	}
	fetcher = publicipapi.NewResilient(fetchers, logger)
	return fetcher, nil
}

func buildFetchers(client *http.Client, parallelResolver common.ParallelResolver,
	unzipper common.Unzipper, ipFetcher common.IPFetcher,
	logger Logger, protonEmail, protonPassword string,
) map[string]common.Fetcher {
	return map[string]common.Fetcher{
		providers.Airvpn:                airvpn.New(client),
		providers.Cyberghost:            cyberghost.New(parallelResolver, logger),
		providers.Expressvpn:            expressvpn.New(unzipper, logger, parallelResolver),
		providers.Fastestvpn:            fastestvpn.New(client, logger, parallelResolver),
		providers.Giganews:              giganews.New(unzipper, logger, parallelResolver),
		providers.HideMyAss:             hidemyass.New(client, logger, parallelResolver),
		providers.Ipvanish:              ipvanish.New(unzipper, logger, parallelResolver),
		providers.Ivpn:                  ivpn.New(client, logger, parallelResolver),
		providers.Mullvad:               mullvad.New(client),
		providers.Nordvpn:               nordvpn.New(client, logger),
		providers.Perfectprivacy:        perfectprivacy.New(unzipper, logger),
		providers.Privado:               privado.New(client, logger),
		providers.PrivateInternetAccess: privateinternetaccess.New(client),
		providers.Privatevpn:            privatevpn.New(unzipper, logger, parallelResolver),
		providers.Protonvpn:             protonvpn.New(client, logger, protonEmail, protonPassword),
		providers.Purevpn:               purevpn.New(ipFetcher, unzipper, logger, parallelResolver),
		providers.SlickVPN:              slickvpn.New(client, logger, parallelResolver),
		providers.Surfshark:             surfshark.New(client, unzipper, logger, parallelResolver),
		providers.Torguard:              torguard.New(unzipper, logger, parallelResolver),
		providers.VPNSecure:             vpnsecure.New(client, logger, parallelResolver),
		providers.VPNUnlimited:          vpnunlimited.New(unzipper, logger, parallelResolver),
		providers.Vyprvpn:               vyprvpn.New(unzipper, logger, parallelResolver),
		providers.Windscribe:            windscribe.New(client, logger),
	}
}

func fetchAndWrite(ctx context.Context, providerName string, fetcher common.Fetcher,
	minServers uint, outputDirPath string,
) error {
	filename := strings.ToLower(strings.ReplaceAll(providerName, " ", "")) + ".json"
	destinationPath := filepath.Join(outputDirPath, filename)
	const permission = 0o644
	file, err := os.OpenFile(destinationPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, permission)
	if err != nil {
		return fmt.Errorf("opening output file: %w", err)
	}

	servers, err := fetcher.FetchServers(ctx, int(minServers)) //nolint:gosec
	if err != nil {
		_ = file.Close()
		_ = os.Remove(destinationPath)
		return fmt.Errorf("fetching servers: %w", err)
	}
	data := models.Servers{
		Version:   fetcher.Version(),
		Timestamp: time.Now().Unix(),
		Servers:   servers,
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(data)
	if err != nil {
		_ = file.Close()
		return fmt.Errorf("encoding servers to JSON: %w", err)
	}

	err = file.Close()
	if err != nil {
		return fmt.Errorf("closing output file: %w", err)
	}

	return nil
}
