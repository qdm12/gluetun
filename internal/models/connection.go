package models

import (
	"net"
)

type Connection struct {
	// Type is the connection type and can be "openvpn" or "wireguard"
	Type string `json:"type"`
	// IP is the VPN server IP address.
	IP net.IP `json:"ip"`
	// Port is the VPN server port.
	Port uint16 `json:"port"`
	// Protocol can be "tcp" or "udp".
	Protocol string `json:"protocol"`
	// Hostname is used for IPVanish, IVPN, Privado
	// and Windscribe for TLS verification.
	Hostname string `json:"hostname"`
	// ServerName is used for PIA for port forwarding
	ServerName string `json:"server_name,omitempty"`
	// PubKey is the public key of the VPN server,
	// used only for Wireguard.
	PubKey string `json:"pubkey"`
}

func (c *Connection) Equal(other Connection) bool {
	return c.IP.Equal(other.IP) && c.Port == other.Port &&
		c.Protocol == other.Protocol && c.Hostname == other.Hostname &&
		c.ServerName == other.ServerName && c.PubKey == other.PubKey
}

// UpdateEmptyWith updates each field of the connection where the
// value is not set using the value given as arguments.
func (c *Connection) UpdateEmptyWith(ip net.IP, port uint16, protocol string) {
	if c.IP == nil {
		c.IP = ip
	}
	if c.Port == 0 {
		c.Port = port
	}
	if c.Protocol == "" {
		c.Protocol = protocol
	}
}
