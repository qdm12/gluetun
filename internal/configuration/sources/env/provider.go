package env

import (
	"fmt"
	"strings"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/constants/vpn"
)

func (s *Source) readProvider(vpnType string) (provider settings.Provider, err error) {
	provider.Name = s.readVPNServiceProvider(vpnType)
	var providerName string
	if provider.Name != nil {
		providerName = *provider.Name
	}

	provider.ServerSelection, err = s.readServerSelection(providerName, vpnType)
	if err != nil {
		return provider, fmt.Errorf("server selection: %w", err)
	}

	provider.PortForwarding, err = s.readPortForward()
	if err != nil {
		return provider, fmt.Errorf("port forwarding: %w", err)
	}

	return provider, nil
}

func (s *Source) readVPNServiceProvider(vpnType string) (vpnProviderPtr *string) {
	_, value := s.getEnvWithRetro("VPN_SERVICE_PROVIDER", "VPNSP")
	if value == "" {
		if vpnType != vpn.Wireguard && getCleanedEnv("OPENVPN_CUSTOM_CONFIG") != "" {
			// retro compatibility
			return stringPtr(providers.Custom)
		}
		return nil
	}

	value = strings.ToLower(value)
	if value == "pia" { // retro compatibility
		return stringPtr(providers.PrivateInternetAccess)
	}

	return stringPtr(value)
}
