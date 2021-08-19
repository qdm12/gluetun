package vpn

import (
	"context"
	"errors"
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/firewall"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/openvpn"
	"github.com/qdm12/gluetun/internal/openvpn/custom"
	"github.com/qdm12/gluetun/internal/provider"
	"github.com/qdm12/golibs/command"
	"github.com/qdm12/golibs/logging"
)

var (
	errBuildConfig = errors.New("failed building configuration")
	errWriteConfig = errors.New("failed writing configuration to file")
	errWriteAuth   = errors.New("failed writing auth to file")
	errFirewall    = errors.New("failed allowing VPN connection through firewall")
)

// setupOpenVPN sets OpenVPN up using the configurators and settings given.
// It returns a serverName for port forwarding (PIA) and an error if it fails.
func setupOpenVPN(ctx context.Context, fw firewall.VPNConnectionSetter,
	openvpnConf openvpn.Interface, providerConf provider.Provider,
	openVPNSettings configuration.OpenVPN, providerSettings configuration.Provider,
	starter command.Starter, logger logging.Logger) (
	runner vpnRunner, serverName string, err error) {
	var connection models.Connection
	var lines []string
	if openVPNSettings.Config == "" {
		connection, err = providerConf.GetConnection(providerSettings.ServerSelection)
		if err == nil {
			lines = providerConf.BuildConf(connection, openVPNSettings)
		}
	} else {
		lines, connection, err = custom.BuildConfig(openVPNSettings)
	}
	if err != nil {
		return nil, "", fmt.Errorf("%w: %s", errBuildConfig, err)
	}

	if err := openvpnConf.WriteConfig(lines); err != nil {
		return nil, "", fmt.Errorf("%w: %s", errWriteConfig, err)
	}

	if openVPNSettings.User != "" {
		err := openvpnConf.WriteAuthFile(openVPNSettings.User, openVPNSettings.Password)
		if err != nil {
			return nil, "", fmt.Errorf("%w: %s", errWriteAuth, err)
		}
	}

	if err := fw.SetVPNConnection(ctx, connection); err != nil {
		return nil, "", fmt.Errorf("%w: %s", errFirewall, err)
	}

	runner = openvpn.NewRunner(openVPNSettings, starter, logger)

	return runner, connection.Hostname, nil
}
