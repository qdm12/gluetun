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

	settings.Endpoint = new(net.UDPAddr)
	settings.Endpoint.IP = make(net.IP, len(connection.IP))
	copy(settings.Endpoint.IP, connection.IP)
	settings.Endpoint.Port = int(connection.Port)

	address := new(net.IPNet)
	address.IP = make(net.IP, len(userSettings.Address.IP))
	copy(address.IP, userSettings.Address.IP)
	address.Mask = make(net.IPMask, len(userSettings.Address.Mask))
	copy(address.Mask, userSettings.Address.Mask)
	settings.Addresses = append(settings.Addresses, address)

	return settings
}
