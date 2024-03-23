package env

import (
	"fmt"
	"strings"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/constants/vpn"
	"github.com/qdm12/gosettings/sources/env"
)

func (s *Source) readProvider(vpnType string) (provider settings.Provider, err error) {
	provider.Name = s.readVPNServiceProvider(vpnType)

	provider.ServerSelection, err = s.readServerSelection(provider.Name, vpnType)
	if err != nil {
		return provider, fmt.Errorf("server selection: %w", err)
	}

	provider.PortForwarding, err = s.readPortForward()
	if err != nil {
		return provider, fmt.Errorf("port forwarding: %w", err)
	}

	return provider, nil
}

func (s *Source) readVPNServiceProvider(vpnType string) (vpnProvider string) {
	vpnProvider = s.env.String("VPN_SERVICE_PROVIDER", env.RetroKeys("VPNSP"))
	if vpnProvider == "" {
		if vpnType != vpn.Wireguard && s.env.Get("OPENVPN_CUSTOM_CONFIG") != nil {
			// retro compatibility
			return providers.Custom
		}
		return ""
	}

	vpnProvider = strings.ToLower(vpnProvider)
	if vpnProvider == "pia" { // retro compatibility
		return providers.PrivateInternetAccess
	}

	return vpnProvider
}
