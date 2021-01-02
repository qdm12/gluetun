package cli

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/settings"
	"github.com/qdm12/gluetun/internal/storage"
	"github.com/qdm12/gluetun/internal/updater"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/os"
)

func (c *cli) Update(args []string, os os.OS) error {
	options := settings.Updater{CLI: true}
	var flushToFile bool
	flagSet := flag.NewFlagSet("update", flag.ExitOnError)
	flagSet.BoolVar(&flushToFile, "file", false, "Write results to /gluetun/servers.json (for end users)")
	flagSet.BoolVar(&options.Stdout, "stdout", false, "Write results to console to modify the program (for maintainers)")
	flagSet.StringVar(&options.DNSAddress, "dns", "1.1.1.1", "DNS resolver address to use")
	flagSet.BoolVar(&options.Cyberghost, "cyberghost", false, "Update Cyberghost servers")
	flagSet.BoolVar(&options.Mullvad, "mullvad", false, "Update Mullvad servers")
	flagSet.BoolVar(&options.Nordvpn, "nordvpn", false, "Update Nordvpn servers")
	flagSet.BoolVar(&options.PIA, "pia", false, "Update Private Internet Access post-summer 2020 servers")
	flagSet.BoolVar(&options.Privado, "privado", false, "Update Privado servers")
	flagSet.BoolVar(&options.Purevpn, "purevpn", false, "Update Purevpn servers")
	flagSet.BoolVar(&options.Surfshark, "surfshark", false, "Update Surfshark servers")
	flagSet.BoolVar(&options.Vyprvpn, "vyprvpn", false, "Update Vyprvpn servers")
	flagSet.BoolVar(&options.Windscribe, "windscribe", false, "Update Windscribe servers")
	if err := flagSet.Parse(args); err != nil {
		return err
	}
	logger, err := logging.NewLogger(logging.ConsoleEncoding, logging.InfoLevel)
	if err != nil {
		return err
	}
	if !flushToFile && !options.Stdout {
		return fmt.Errorf("at least one of -file or -stdout must be specified")
	}
	ctx := context.Background()
	const clientTimeout = 10 * time.Second
	httpClient := &http.Client{Timeout: clientTimeout}
	storage := storage.New(logger, os, constants.ServersData)
	currentServers, err := storage.SyncServers(constants.GetAllServers())
	if err != nil {
		return fmt.Errorf("cannot update servers: %w", err)
	}
	updater := updater.New(options, httpClient, currentServers, logger)
	allServers, err := updater.UpdateServers(ctx)
	if err != nil {
		return err
	}
	if flushToFile {
		if err := storage.FlushToFile(allServers); err != nil {
			return fmt.Errorf("cannot update servers: %w", err)
		}
	}

	return nil
}
