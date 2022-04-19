package utils

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/constants/vpn"
)

func getProtocol(selection settings.ServerSelection) (protocol string) {
	if selection.VPN == vpn.OpenVPN && *selection.OpenVPN.TCP {
		return constants.TCP
	}
	return constants.UDP
}

func filterByProtocol(selection settings.ServerSelection,
	serverTCP, serverUDP bool) (filtered bool) {
	switch selection.VPN {
	case vpn.Wireguard:
		return !serverUDP
	default: // OpenVPN
		wantTCP := *selection.OpenVPN.TCP
		wantUDP := !wantTCP
		return (wantTCP && !serverTCP) || (wantUDP && !serverUDP)
	}
}
