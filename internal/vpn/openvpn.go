package vpn

import (
	"context"
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/openvpn"
	"github.com/qdm12/gluetun/internal/provider"
	"github.com/qdm12/golibs/command"
)

// setupOpenVPN sets OpenVPN up using the configurators and settings given.
// It returns a serverName for port forwarding (PIA) and an error if it fails.
func setupOpenVPN(ctx context.Context, fw Firewall,
	openvpnConf OpenVPN, providerConf provider.Provider,
	settings settings.VPN, ipv6Supported bool, starter command.Starter,
	logger openvpn.Logger) (runner *openvpn.Runner, serverName string, err error) {
	connection, err := providerConf.GetConnection(settings.Provider.ServerSelection, ipv6Supported)
	if err != nil {
		return nil, "", fmt.Errorf("failed finding a valid server connection: %w", err)
	}

	lines := providerConf.OpenVPNConfig(connection, settings.OpenVPN, ipv6Supported)

	if err := openvpnConf.WriteConfig(lines); err != nil {
		return nil, "", fmt.Errorf("failed writing configuration to file: %w", err)
	}

	if *settings.OpenVPN.User != "" {
		err := openvpnConf.WriteAuthFile(*settings.OpenVPN.User, *settings.OpenVPN.Password)
		if err != nil {
			return nil, "", fmt.Errorf("failed writing auth to file: %w", err)
		}
	}

	if *settings.OpenVPN.KeyPassphrase != "" {
		err := openvpnConf.WriteAskPassFile(*settings.OpenVPN.KeyPassphrase)
		if err != nil {
			return nil, "", fmt.Errorf("writing askpass file: %w", err)
		}
	}

	if err := fw.SetVPNConnection(ctx, connection, settings.OpenVPN.Interface); err != nil {
		return nil, "", fmt.Errorf("failed allowing VPN connection through firewall: %w", err)
	}

	runner = openvpn.NewRunner(settings.OpenVPN, starter, logger)

	return runner, connection.ServerName, nil
}
