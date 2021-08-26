package cli

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/storage"
	"github.com/qdm12/gluetun/internal/updater"
	"github.com/qdm12/golibs/logging"
)

var (
	ErrModeUnspecified         = errors.New("at least one of -enduser or -maintainers must be specified")
	ErrNewStorage              = errors.New("cannot create storage")
	ErrUpdateServerInformation = errors.New("cannot update server information")
	ErrWriteToFile             = errors.New("cannot write updated information to file")
)

type Updater interface {
	Update(ctx context.Context, args []string, logger logging.Logger) error
}

func (c *CLI) Update(ctx context.Context, args []string, logger logging.Logger) error {
	options := configuration.Updater{CLI: true}
	var endUserMode, maintainerMode, updateAll bool
	flagSet := flag.NewFlagSet("update", flag.ExitOnError)
	flagSet.BoolVar(&endUserMode, "enduser", false, "Write results to /gluetun/servers.json (for end users)")
	flagSet.BoolVar(&maintainerMode, "maintainer", false,
		"Write results to ./internal/storage/servers.json to modify the program (for maintainers)")
	flagSet.StringVar(&options.DNSAddress, "dns", "8.8.8.8", "DNS resolver address to use")
	flagSet.BoolVar(&updateAll, "all", false, "Update servers for all VPN providers")
	flagSet.BoolVar(&options.Cyberghost, "cyberghost", false, "Update Cyberghost servers")
	flagSet.BoolVar(&options.Fastestvpn, "fastestvpn", false, "Update FastestVPN servers")
	flagSet.BoolVar(&options.HideMyAss, "hidemyass", false, "Update HideMyAss servers")
	flagSet.BoolVar(&options.Ipvanish, "ipvanish", false, "Update IpVanish servers")
	flagSet.BoolVar(&options.Ivpn, "ivpn", false, "Update IVPN servers")
	flagSet.BoolVar(&options.Mullvad, "mullvad", false, "Update Mullvad servers")
	flagSet.BoolVar(&options.Nordvpn, "nordvpn", false, "Update Nordvpn servers")
	flagSet.BoolVar(&options.PIA, "pia", false, "Update Private Internet Access post-summer 2020 servers")
	flagSet.BoolVar(&options.Privado, "privado", false, "Update Privado servers")
	flagSet.BoolVar(&options.Privatevpn, "privatevpn", false, "Update Private VPN servers")
	flagSet.BoolVar(&options.Protonvpn, "protonvpn", false, "Update Protonvpn servers")
	flagSet.BoolVar(&options.Purevpn, "purevpn", false, "Update Purevpn servers")
	flagSet.BoolVar(&options.Surfshark, "surfshark", false, "Update Surfshark servers")
	flagSet.BoolVar(&options.Torguard, "torguard", false, "Update Torguard servers")
	flagSet.BoolVar(&options.VPNUnlimited, "vpnunlimited", false, "Update VPN Unlimited servers")
	flagSet.BoolVar(&options.Vyprvpn, "vyprvpn", false, "Update Vyprvpn servers")
	flagSet.BoolVar(&options.Wevpn, "wevpn", false, "Update WeVPN servers")
	flagSet.BoolVar(&options.Windscribe, "windscribe", false, "Update Windscribe servers")
	if err := flagSet.Parse(args); err != nil {
		return err
	}
	if !endUserMode && !maintainerMode {
		return ErrModeUnspecified
	}

	if updateAll {
		options.EnableAll()
	}

	const clientTimeout = 10 * time.Second
	httpClient := &http.Client{Timeout: clientTimeout}

	storage, err := storage.New(logger, constants.ServersData)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrNewStorage, err)
	}
	currentServers := storage.GetServers()

	updater := updater.New(options, httpClient, currentServers, logger)
	allServers, err := updater.UpdateServers(ctx)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUpdateServerInformation, err)
	}

	if endUserMode {
		if err := storage.FlushToFile(allServers); err != nil {
			return fmt.Errorf("%w: %s", ErrWriteToFile, err)
		}
	}

	if maintainerMode {
		if err := writeToEmbeddedJSON(c.repoServersPath, allServers); err != nil {
			return fmt.Errorf("%w: %s", ErrWriteToFile, err)
		}
	}

	return nil
}

func writeToEmbeddedJSON(repoServersPath string,
	allServers models.AllServers) error {
	const perms = 0600
	f, err := os.OpenFile(repoServersPath,
		os.O_TRUNC|os.O_WRONLY|os.O_CREATE, perms)
	if err != nil {
		return err
	}

	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	return encoder.Encode(allServers)
}
