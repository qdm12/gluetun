package utils

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants/vpn"
)

func getPort(selection settings.ServerSelection,
	defaultOpenVPNTCP, defaultOpenVPNUDP, defaultWireguard uint16) (port uint16) {
	switch selection.VPN {
	case vpn.Wireguard:
		customPort := *selection.Wireguard.EndpointPort
		if customPort > 0 {
			return customPort
		}
		checkDefined("Wireguard", defaultWireguard)
		return defaultWireguard
	default: // OpenVPN
		customPort := *selection.OpenVPN.CustomPort
		if customPort > 0 {
			return customPort
		}
		if *selection.OpenVPN.TCP {
			checkDefined("OpenVPN TCP", defaultOpenVPNTCP)
			return defaultOpenVPNTCP
		}

		checkDefined("OpenVPN UDP", defaultOpenVPNUDP)
		return defaultOpenVPNUDP
	}
}

func checkDefined(portName string, port uint16) {
	if port > 0 {
		return
	}
	message := fmt.Sprintf("no default %s port is defined!", portName)
	panic(message)
}
