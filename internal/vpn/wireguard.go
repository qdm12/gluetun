package vpn

import (
	"context"
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/provider"
	"github.com/qdm12/gluetun/internal/provider/utils"
	"github.com/qdm12/gluetun/internal/wireguard"
)

// setupWireguard sets Wireguard up using the configurators and settings given.
// It returns a serverName for port forwarding (PIA) and an error if it fails.
func setupWireguard(ctx context.Context, netlinker NetLinker,
	fw Firewall, providerConf provider.Provider,
	settings settings.VPN, ipv6Supported bool, logger wireguard.Logger) (
	wireguarder *wireguard.Wireguard, serverName string, err error) {
	connection, err := providerConf.GetConnection(settings.Provider.ServerSelection, ipv6Supported)
	if err != nil {
		return nil, "", fmt.Errorf("failed finding a VPN server: %w", err)
	}

	wireguardSettings := utils.BuildWireguardSettings(connection, settings.Wireguard, ipv6Supported)

	logger.Debug("Wireguard server public key: " + wireguardSettings.PublicKey)
	logger.Debug("Wireguard client private key: " + wireguardSettings.PrivateKey)
	logger.Debug("Wireguard pre-shared key: " + wireguardSettings.PreSharedKey)

	wireguarder, err = wireguard.New(wireguardSettings, netlinker, logger)
	if err != nil {
		return nil, "", fmt.Errorf("failed creating Wireguard: %w", err)
	}

	err = fw.SetVPNConnection(ctx, connection, settings.Wireguard.Interface)
	if err != nil {
		return nil, "", fmt.Errorf("failed setting firewall: %w", err)
	}

	return wireguarder, connection.ServerName, nil
}
