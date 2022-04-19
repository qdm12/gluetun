package models

import (
	"net"

	"github.com/qdm12/gluetun/internal/constants/vpn"
)

type Server struct {
	VPN string `json:"vpn,omitempty"`
	// Surfshark: country is also used for multi-hop
	Country     string   `json:"country,omitempty"`
	Region      string   `json:"region,omitempty"`
	City        string   `json:"city,omitempty"`
	ISP         string   `json:"isp,omitempty"`
	Owned       bool     `json:"owned,omitempty"`
	Number      uint16   `json:"number,omitempty"`
	ServerName  string   `json:"server_name,omitempty"`
	Hostname    string   `json:"hostname,omitempty"`
	TCP         bool     `json:"tcp,omitempty"`
	UDP         bool     `json:"udp,omitempty"`
	OvpnX509    string   `json:"x509,omitempty"`
	RetroLoc    string   `json:"retroloc,omitempty"` // TODO remove in v4
	MultiHop    bool     `json:"multihop,omitempty"`
	WgPubKey    string   `json:"wgpubkey,omitempty"`
	Free        bool     `json:"free,omitempty"`
	Stream      bool     `json:"stream,omitempty"`
	PortForward bool     `json:"port_forward,omitempty"`
	IPs         []net.IP `json:"ips,omitempty"`
}

func (s *Server) setDefaults() {
	// TODO v4 precise these in servers.json rather than here
	if s.VPN == "" {
		// If the VPN protocol isn't specified, assume it is OpenVPN.
		s.VPN = vpn.OpenVPN
	}

	if !s.UDP && !s.TCP {
		// If UDP and TCP are not precised:
		// For OpenVPN, assume TCP and UDP are supported
		// For Wireguard, UDP must be supported.
		s.UDP = true
		if s.VPN == vpn.OpenVPN {
			s.TCP = true
		}
	}
}
