package utils

import (
	"net"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/wireguard"
)

func BuildWireguardSettings(connection models.Connection,
	userSettings configuration.Wireguard) (settings wireguard.Settings) {
	settings.PrivateKey = userSettings.PrivateKey
	settings.PublicKey = connection.PubKey
	settings.PreSharedKey = userSettings.PreSharedKey
	settings.InterfaceName = userSettings.Interface

	const routePriority = 101 // 100 is to receive external connections
	settings.RulePriority = routePriority

	settings.Endpoint = new(net.UDPAddr)
	settings.Endpoint.IP = make(net.IP, len(connection.IP))
	copy(settings.Endpoint.IP, connection.IP)
	settings.Endpoint.Port = int(connection.Port)

	for _, address := range settings.Addresses {
		addressCopy := new(net.IPNet)
		addressCopy.IP = make(net.IP, len(address.IP))
		copy(addressCopy.IP, address.IP)
		addressCopy.Mask = make(net.IPMask, len(address.Mask))
		copy(addressCopy.Mask, address.Mask)
		settings.Addresses = append(settings.Addresses, addressCopy)
	}

	return settings
}
