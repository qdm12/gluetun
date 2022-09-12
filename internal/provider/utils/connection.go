package utils

import (
	"fmt"
	"math/rand"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants/vpn"
	"github.com/qdm12/gluetun/internal/models"
)

type ConnectionDefaults struct {
	OpenVPNTCPPort uint16
	OpenVPNUDPPort uint16
	WireguardPort  uint16
}

func NewConnectionDefaults(openvpnTCPPort, openvpnUDPPort,
	wireguardPort uint16) ConnectionDefaults {
	return ConnectionDefaults{
		OpenVPNTCPPort: openvpnTCPPort,
		OpenVPNUDPPort: openvpnUDPPort,
		WireguardPort:  wireguardPort,
	}
}

type Storage interface {
	FilterServers(provider string, selection settings.ServerSelection) (
		servers []models.Server, err error)
}

func GetConnection(provider string,
	storage Storage,
	selection settings.ServerSelection,
	defaults ConnectionDefaults,
	ipv6Supported bool,
	randSource rand.Source) (
	connection models.Connection, err error) {
	servers, err := storage.FilterServers(provider, selection)
	if err != nil {
		return connection, fmt.Errorf("cannot filter servers: %w", err)
	}

	protocol := getProtocol(selection)
	port := getPort(selection, defaults.OpenVPNTCPPort,
		defaults.OpenVPNUDPPort, defaults.WireguardPort)

	connections := make([]models.Connection, 0, len(servers))
	for _, server := range servers {
		for _, ip := range server.IPs {
			if !ipv6Supported && ip.To4() == nil {
				continue
			}

			hostname := server.Hostname
			if selection.VPN == vpn.OpenVPN && server.OvpnX509 != "" {
				// For Windscribe where hostname and
				// OpenVPN x509 are not the same.
				hostname = server.OvpnX509
			}

			connection := models.Connection{
				Type:       selection.VPN,
				IP:         ip,
				Port:       port,
				Protocol:   protocol,
				Hostname:   hostname,
				ServerName: server.ServerName,
				PubKey:     server.WgPubKey, // Wireguard
			}
			connections = append(connections, connection)
		}
	}

	return pickConnection(connections, selection, randSource)
}
