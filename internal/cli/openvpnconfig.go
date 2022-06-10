package cli

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/qdm12/gluetun/internal/configuration/sources"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/provider"
	"github.com/qdm12/gluetun/internal/storage"
	"github.com/qdm12/gluetun/internal/updater/unzip"
)

type OpenvpnConfigMaker interface {
	OpenvpnConfig(logger OpenvpnConfigLogger, source sources.Source) error
}

type OpenvpnConfigLogger interface {
	Info(s string)
	Warn(s string)
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
	unzipper := (unzip.Unzipper)(nil)
	client := (*http.Client)(nil)
	warner := (Warner)(nil)

	providers := provider.NewProviders(storage, time.Now, warner, client, unzipper)
	providerConf := providers.Get(*allSettings.VPN.Provider.Name)
	connection, err := providerConf.GetConnection(allSettings.VPN.Provider.ServerSelection)
	if err != nil {
		return err
	}

	lines := providerConf.OpenVPNConfig(connection, allSettings.VPN.OpenVPN)

	fmt.Println(strings.Join(lines, "\n"))
	return nil
}
