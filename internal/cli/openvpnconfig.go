package cli

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/qdm12/gluetun/internal/configuration/sources"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/provider"
	publicipmodels "github.com/qdm12/gluetun/internal/publicip/models"
	"github.com/qdm12/gluetun/internal/storage"
	"github.com/qdm12/gluetun/internal/updater/resolver"
)

type OpenvpnConfigLogger interface {
	Info(s string)
	Warn(s string)
}

type Unzipper interface {
	FetchAndExtract(ctx context.Context, url string) (
		contents map[string][]byte, err error)
}

type ParallelResolver interface {
	Resolve(ctx context.Context, settings resolver.ParallelSettings) (
		hostToIPs map[string][]net.IP, warnings []string, err error)
}

type IPFetcher interface {
	FetchMultiInfo(ctx context.Context, ips []net.IP) (data []publicipmodels.IPInfoData, err error)
}

func (c *CLI) OpenvpnConfig(logger OpenvpnConfigLogger, source sources.Source) error {
	storage, err := storage.New(logger, constants.ServersData)
	if err != nil {
		return err
	}

	allSettings, err := source.Read()
	if err != nil {
		return err
	}

	if err = allSettings.Validate(storage); err != nil {
		return err
	}

	// Unused by this CLI command
	unzipper := (Unzipper)(nil)
	client := (*http.Client)(nil)
	warner := (Warner)(nil)
	parallelResolver := (ParallelResolver)(nil)
	ipFetcher := (IPFetcher)(nil)

	providers := provider.NewProviders(storage, time.Now, warner, client,
		unzipper, parallelResolver, ipFetcher)
	providerConf := providers.Get(*allSettings.VPN.Provider.Name)
	connection, err := providerConf.GetConnection(allSettings.VPN.Provider.ServerSelection)
	if err != nil {
		return err
	}

	lines := providerConf.OpenVPNConfig(connection, allSettings.VPN.OpenVPN)

	fmt.Println(strings.Join(lines, "\n"))
	return nil
}
