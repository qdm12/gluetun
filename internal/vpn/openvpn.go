package vpn

import (
	"context"
	"errors"
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/firewall"
	"github.com/qdm12/gluetun/internal/openvpn"
	"github.com/qdm12/gluetun/internal/provider"
	"github.com/qdm12/golibs/command"
)

var (
	errServerConn  = errors.New("failed finding a valid server connection")
	errBuildConfig = errors.New("failed building configuration")
	errWriteConfig = errors.New("failed writing configuration to file")
	errWriteAuth   = errors.New("failed writing auth to file")
	errFirewall    = errors.New("failed allowing VPN connection through firewall")
)

// setupOpenVPN sets OpenVPN up using the configurators and settings given.
// It returns a serverName for port forwarding (PIA) and an error if it fails.
func setupOpenVPN(ctx context.Context, fw firewall.VPNConnectionSetter,
	openvpnConf openvpn.Interface, providerConf provider.Provider,
	settings settings.VPN, starter command.Starter, logger openvpn.Logger) (
	runner vpnRunner, serverName string, err error) {
	connection, err := providerConf.GetConnection(settings.Provider.ServerSelection)
	if err != nil {
		return nil, "", fmt.Errorf("%w: %s", errServerConn, err)
	}

	lines, err := providerConf.BuildConf(connection, settings.OpenVPN)
	if err != nil {
		return nil, "", fmt.Errorf("%w: %s", errBuildConfig, err)
	}

	if err := openvpnConf.WriteConfig(lines); err != nil {
		return nil, "", fmt.Errorf("%w: %s", errWriteConfig, err)
	}

	if settings.OpenVPN.User != "" {
		err := openvpnConf.WriteAuthFile(settings.OpenVPN.User, settings.OpenVPN.Password)
		if err != nil {
			return nil, "", fmt.Errorf("%w: %s", errWriteAuth, err)
		}
	}

	if err := fw.SetVPNConnection(ctx, connection, settings.OpenVPN.Interface); err != nil {
		return nil, "", fmt.Errorf("%w: %s", errFirewall, err)
	}

	runner = openvpn.NewRunner(settings.OpenVPN, starter, logger)

	return runner, connection.Hostname, nil
}
