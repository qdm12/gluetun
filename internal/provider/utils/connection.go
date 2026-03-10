package utils

import (
	"fmt"
	"math/rand"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/constants/vpn"
	"github.com/qdm12/gluetun/internal/models"
)

type ConnectionDefaults struct {
	OpenVPNTCPPort uint16
	OpenVPNUDPPort uint16
	WireguardPort  uint16
}

func NewConnectionDefaults(openvpnTCPPort, openvpnUDPPort,
	wireguardPort uint16,
) ConnectionDefaults {
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
	connection models.Connection, err error,
) {
	servers, err := storage.FilterServers(provider, selection)
	if err != nil {
		return connection, fmt.Errorf("filtering servers: %w", err)
	}

	protocol := getProtocol(selection)
	port := getPort(selection, defaults.OpenVPNTCPPort,
		defaults.OpenVPNUDPPort, defaults.WireguardPort)

	connections := make([]models.Connection, 0, len(servers))
	for _, server := range servers {
		for _, ip := range server.IPs {
			if !ipv6Supported && ip.Is6() {
				continue
			}

			hostname := server.Hostname
			if selection.VPN == vpn.OpenVPN && server.OvpnX509 != "" {
				// For Windscribe where hostname and
				// OpenVPN x509 are not the same.
				hostname = server.OvpnX509
			}

			portForServer := port
			customOpenVPNPortSet := selection.OpenVPN.CustomPort != nil &&
				*selection.OpenVPN.CustomPort != 0
			if !customOpenVPNPortSet && selection.VPN == vpn.OpenVPN {
				portForServer = getPortForServer(server, protocol,
					defaults.OpenVPNTCPPort, defaults.OpenVPNUDPPort)
			}

			connection := models.Connection{
				Type:        selection.VPN,
				IP:          ip,
				Port:        portForServer,
				Protocol:    protocol,
				Hostname:    hostname,
				ServerName:  server.ServerName,
				PortForward: server.PortForward,
				PubKey:      server.WgPubKey, // Wireguard
			}
			connections = append(connections, connection)
		}
	}

	return pickConnection(connections, selection, randSource)
}

func getPortForServer(server models.Server, protocol string, defaultTCPPort, defaultUDPPort uint16) (port uint16) {
	switch protocol {
	case constants.TCP:
		ports := make([]uint16, 0, len(server.TCPPorts)+3)
		ports = append(ports, server.TCPPorts...)
		ports = append(ports, defaultTCPPort, 443, 1194)
		return firstNonZeroPort(ports)
	case constants.UDP:
		ports := make([]uint16, 0, len(server.UDPPorts)+3)
		ports = append(ports, server.UDPPorts...)
		ports = append(ports, defaultUDPPort, 1194, 53)
		return firstNonZeroPort(ports)
	default:
		return 0
	}
}

func firstNonZeroPort(ports []uint16) (port uint16) {
	for _, port := range ports {
		if port != 0 {
			return port
		}
	}
	return 0
}
