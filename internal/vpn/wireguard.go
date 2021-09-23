package vpn

import (
	"context"
	"errors"
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/firewall"
	"github.com/qdm12/gluetun/internal/netlink"
	"github.com/qdm12/gluetun/internal/provider"
	"github.com/qdm12/gluetun/internal/provider/utils"
	"github.com/qdm12/gluetun/internal/wireguard"
)

var (
	errGetServer       = errors.New("failed finding a VPN server")
	errCreateWireguard = errors.New("failed creating Wireguard")
)

// setupWireguard sets Wireguard up using the configurators and settings given.
// It returns a serverName for port forwarding (PIA) and an error if it fails.
func setupWireguard(ctx context.Context, netlinker netlink.NetLinker,
	fw firewall.VPNConnectionSetter, providerConf provider.Provider,
	settings configuration.VPN, logger wireguard.Logger) (
	wireguarder wireguard.Wireguarder, serverName string, err error) {
	connection, err := providerConf.GetConnection(settings.Provider.ServerSelection)
	if err != nil {
		return nil, "", fmt.Errorf("%w: %s", errGetServer, err)
	}

	wireguardSettings := utils.BuildWireguardSettings(connection, settings.Wireguard)

	logger.Debug("Wireguard server public key: " + wireguardSettings.PublicKey)
	logger.Debug("Wireguard client private key: " + wireguardSettings.PrivateKey)
	logger.Debug("Wireguard pre-shared key: " + wireguardSettings.PreSharedKey)

	wireguarder, err = wireguard.New(wireguardSettings, netlinker, logger)
	if err != nil {
		return nil, "", fmt.Errorf("%w: %s", errCreateWireguard, err)
	}

	err = fw.SetVPNConnection(ctx, connection, settings.Wireguard.Interface)
	if err != nil {
		return nil, "", fmt.Errorf("%w: %s", errFirewall, err)
	}

	return wireguarder, connection.Hostname, nil
}
