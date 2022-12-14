package utils

import (
	"net"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/wireguard"
)

func BuildWireguardSettings(connection models.Connection,
	userSettings settings.Wireguard, ipv6Supported bool) (settings wireguard.Settings) {
	settings.PrivateKey = *userSettings.PrivateKey
	settings.PublicKey = connection.PubKey
	settings.PreSharedKey = *userSettings.PreSharedKey
	settings.InterfaceName = userSettings.Interface
	settings.Implementation = userSettings.Implementation
	settings.IPv6 = &ipv6Supported

	const rulePriority = 101 // 100 is to receive external connections
	settings.RulePriority = rulePriority

	settings.Endpoint = new(net.UDPAddr)
	settings.Endpoint.IP = make(net.IP, len(connection.IP))
	copy(settings.Endpoint.IP, connection.IP)
	settings.Endpoint.Port = int(connection.Port)

	settings.Addresses = make([]*net.IPNet, 0, len(userSettings.Addresses))
	for _, address := range userSettings.Addresses {
		ipv6Address := address.IP.To4() == nil
		if !ipv6Supported && ipv6Address {
			continue
		}
		addressCopy := new(net.IPNet)
		addressCopy.IP = make(net.IP, len(address.IP))
		copy(addressCopy.IP, address.IP)
		addressCopy.Mask = make(net.IPMask, len(address.Mask))
		copy(addressCopy.Mask, address.Mask)
		settings.Addresses = append(settings.Addresses, addressCopy)
	}

	return settings
}
