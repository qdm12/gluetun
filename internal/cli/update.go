package cli

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/storage"
	"github.com/qdm12/gluetun/internal/updater"
)

var (
	ErrModeUnspecified     = errors.New("at least one of -enduser or -maintainer must be specified")
	ErrDNSAddress          = errors.New("DNS address is not valid")
	ErrNoProviderSpecified = errors.New("no provider was specified")
)

type Updater interface {
	Update(ctx context.Context, args []string, logger UpdaterLogger) error
}

type UpdaterLogger interface {
	Info(s string)
	Warn(s string)
	Error(s string)
}

func boolPtr(b bool) *bool { return &b }

func (c *CLI) Update(ctx context.Context, args []string, logger UpdaterLogger) error {
	options := settings.Updater{CLI: boolPtr(true)}
	var endUserMode, maintainerMode, updateAll bool
	var dnsAddress, csvProviders string
	flagSet := flag.NewFlagSet("update", flag.ExitOnError)
	flagSet.BoolVar(&endUserMode, "enduser", false, "Write results to /gluetun/servers.json (for end users)")
	flagSet.BoolVar(&maintainerMode, "maintainer", false,
		"Write results to ./internal/storage/servers.json to modify the program (for maintainers)")
	flagSet.StringVar(&dnsAddress, "dns", "8.8.8.8", "DNS resolver address to use")
	flagSet.BoolVar(&updateAll, "all", false, "Update servers for all VPN providers")
	flagSet.StringVar(&csvProviders, "providers", "", "CSV string of VPN providers to update server data for")
	if err := flagSet.Parse(args); err != nil {
		return err
	}

	if !endUserMode && !maintainerMode {
		return ErrModeUnspecified
	}

	options.DNSAddress = net.ParseIP(dnsAddress)
	if options.DNSAddress == nil {
		return fmt.Errorf("%w: %s", ErrDNSAddress, dnsAddress)
	}

	if updateAll {
		for _, provider := range constants.AllProviders() {
			if provider == constants.Custom {
				continue
			}
			options.Providers = append(options.Providers, provider)
		}
	} else {
		if csvProviders == "" {
			return ErrNoProviderSpecified
		}
		options.Providers = strings.Split(csvProviders, ",")
	}

	options.SetDefaults(options.Providers[0])

	err := options.Validate()
	if err != nil {
		return fmt.Errorf("options validation failed: %w", err)
	}

	const clientTimeout = 10 * time.Second
	httpClient := &http.Client{Timeout: clientTimeout}

	storage, err := storage.New(logger, constants.ServersData)
	if err != nil {
		return fmt.Errorf("cannot create servers storage: %w", err)
	}
	currentServers := storage.GetServers()

	updater := updater.New(options, httpClient, currentServers, logger)
	allServers, err := updater.UpdateServers(ctx)
	if err != nil {
		return fmt.Errorf("cannot update server information: %w", err)
	}

	if endUserMode {
		if err := storage.FlushToFile(allServers); err != nil {
			return fmt.Errorf("cannot write updated information to file: %w", err)
		}
	}

	if maintainerMode {
		if err := writeToEmbeddedJSON(c.repoServersPath, allServers); err != nil {
			return fmt.Errorf("cannot write updated information to file: %w", err)
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
