package utils

import (
	"fmt"
	"math/rand"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants/vpn"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gosettings/reader"
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

type VPNSettings struct {
	IPv6Server *bool
}

// Read method to populate the VPNSettings from the reader
func (v *VPNSettings) read(reader *reader.Reader) (err error) {
	v.IPv6Server, err = reader.BoolPtr("VPN_IPV6_SERVER")
	return err
}

func GetConnection(provider string,
	storage Storage,
	selection settings.ServerSelection,
	defaults ConnectionDefaults,
	ipv6Supported bool,
	randSource rand.Source,
	reader *reader.Reader) (
	connection models.Connection, err error,
) {
	// Create an instance of VPNSettings and read settings
	var vpnSettings VPNSettings
	if err := vpnSettings.read(reader); err != nil {
		return connection, fmt.Errorf("reading VPN settings: %w", err)
	}

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
			// Skip IPv6 if unsupported or if VPN_IPV6_SERVER is false
			if !ipv6Supported || (vpnSettings.IPv6Server != nil && !*vpnSettings.IPv6Server && ip.Is6()) {
				continue
			}

			hostname := server.Hostname
			if selection.VPN == vpn.OpenVPN && server.OvpnX509 != "" {
				hostname = server.OvpnX509
			}

			connection := models.Connection{
				Type:        selection.VPN,
				IP:          ip,
				Port:        port,
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
