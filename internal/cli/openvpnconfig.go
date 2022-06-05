package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/qdm12/gluetun/internal/configuration/sources"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/provider"
	"github.com/qdm12/gluetun/internal/storage"
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

	providerConf := provider.New(*allSettings.VPN.Provider.Name, storage, time.Now)
	connection, err := providerConf.GetConnection(allSettings.VPN.Provider.ServerSelection)
	if err != nil {
		return err
	}

	lines := providerConf.OpenVPNConfig(connection, allSettings.VPN.OpenVPN)

	fmt.Println(strings.Join(lines, "\n"))
	return nil
}
