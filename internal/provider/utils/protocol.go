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
