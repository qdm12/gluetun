package utils

import (
	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
)

func GetProtocol(selection configuration.ServerSelection) (protocol string) {
	if selection.VPN == constants.OpenVPN && selection.OpenVPN.TCP {
		return constants.TCP
	}
	return constants.UDP
}

func FilterByProtocol(selection configuration.ServerSelection,
	serverTCP, serverUDP bool) (filtered bool) {
	switch selection.VPN {
	case constants.Wireguard:
		return !serverUDP
	default: // OpenVPN
		wantTCP := selection.OpenVPN.TCP
		wantUDP := !wantTCP
		return (wantTCP && !serverTCP) || (wantUDP && !serverUDP)
	}
}
