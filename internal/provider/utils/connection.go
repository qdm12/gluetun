package utils

import (
	"fmt"
	"math/rand"

	"github.com/qdm12/gluetun/internal/configuration/settings"
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

func GetConnection(servers []models.Server,
	selection settings.ServerSelection,
	defaults ConnectionDefaults,
	randSource rand.Source) (
	connection models.Connection, err error) {
	servers, err = FilterServers(servers, selection)
	if err != nil {
		return connection, fmt.Errorf("cannot filter servers: %w", err)
	}

	protocol := getProtocol(selection)
	port := GetPort(selection, defaults.OpenVPNTCPPort,
		defaults.OpenVPNUDPPort, defaults.WireguardPort)

	connections := make([]models.Connection, 0, len(servers))
	for _, server := range servers {
		for _, ip := range server.IPs {
			if ip.To4() == nil {
				// do not use IPv6 connections for now
				continue
			}
			connection := models.Connection{
				Type:     selection.VPN,
				IP:       ip,
				Port:     port,
				Protocol: protocol,
				Hostname: server.Hostname,
				PubKey:   server.WgPubKey, // Wireguard
			}
			connections = append(connections, connection)
		}
	}

	return PickConnection(connections, selection, randSource)
}
