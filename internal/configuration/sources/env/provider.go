package env

import (
	"fmt"
	"os"
	"strings"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
)

func (r *Reader) readProvider(vpnType string) (provider settings.Provider, err error) {
	provider.Name = r.readVPNServiceProvider(vpnType)
	var providerName string
	if provider.Name != nil {
		providerName = *provider.Name
	}

	provider.ServerSelection, err = r.readServerSelection(providerName, vpnType)
	if err != nil {
		return provider, fmt.Errorf("cannot read server selection settings: %w", err)
	}

	provider.PortForwarding, err = r.readPortForward()
	if err != nil {
		return provider, fmt.Errorf("cannot read port forwarding settings: %w", err)
	}

	return provider, nil
}

func (r *Reader) readVPNServiceProvider(vpnType string) (vpnProviderPtr *string) {
	_, s := r.getEnvWithRetro("VPN_SERVICE_PROVIDER", "VPNSP")
	s = strings.ToLower(s)
	switch {
	case vpnType != constants.Wireguard &&
		os.Getenv("OPENVPN_CUSTOM_CONFIG") != "": // retro compatibility
		return stringPtr(constants.Custom)
	case s == "":
		return nil
	case s == "pia": // retro compatibility
		return stringPtr(constants.PrivateInternetAccess)
	}
	return stringPtr(s)
}
