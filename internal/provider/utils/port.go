package utils

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants/vpn"
)

func GetPort(selection settings.ServerSelection,
	defaultOpenVPNTCP, defaultOpenVPNUDP, defaultWireguard uint16) (port uint16) {
	switch selection.VPN {
	case vpn.Wireguard:
		customPort := *selection.Wireguard.EndpointPort
		if customPort > 0 {
			return customPort
		}
		return defaultWireguard
	default: // OpenVPN
		customPort := *selection.OpenVPN.CustomPort
		if customPort > 0 {
			return customPort
		}
		if *selection.OpenVPN.TCP {
			return defaultOpenVPNTCP
		}
		return defaultOpenVPNUDP
	}
}
