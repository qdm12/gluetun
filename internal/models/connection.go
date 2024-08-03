package models

import (
	"net/netip"
)

type Connection struct {
	// Type is the connection type and can be "openvpn" or "wireguard"
	Type string `json:"type"`
	// IP is the VPN server IP address.
	IP netip.Addr `json:"ip"`
	// Port is the VPN server port.
	Port uint16 `json:"port"`
	// Protocol can be "tcp" or "udp".
	Protocol string `json:"protocol"`
	// Hostname is used for IPVanish, IVPN, Privado
	// and Windscribe for TLS verification.
	Hostname string `json:"hostname"`
	// PubKey is the public key of the VPN server,
	// used only for Wireguard.
	PubKey string `json:"pubkey"`
	// ServerName is used for PIA for port forwarding
	ServerName string `json:"server_name,omitempty"`
	// PortForward is used for PIA and ProtonVPN for port forwarding
	PortForward bool `json:"port_forward"`
}

func (c *Connection) Equal(other Connection) bool {
	return c.IP.Compare(other.IP) == 0 && c.Port == other.Port &&
		c.Protocol == other.Protocol && c.Hostname == other.Hostname &&
		c.PubKey == other.PubKey && c.ServerName == other.ServerName &&
		c.PortForward == other.PortForward
}

// UpdateEmptyWith updates each field of the connection where the
// value is not set using the value given as arguments.
func (c *Connection) UpdateEmptyWith(ip netip.Addr, port uint16, protocol string) {
	if !c.IP.IsValid() {
		c.IP = ip
	}
	if c.Port == 0 {
		c.Port = port
	}
	if c.Protocol == "" {
		c.Protocol = protocol
	}
}
