package models

import (
	"fmt"
	"net"
	"reflect"
	"strings"
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
	Keep        bool     `json:"keep,omitempty"`
	IPs         []net.IP `json:"ips,omitempty"`
}

func (s *Server) Equal(other Server) (equal bool) {
	if !ipsAreEqual(s.IPs, other.IPs) {
		return false
	}

	serverCopy := *s
	serverCopy.IPs = nil
	other.IPs = nil
	return reflect.DeepEqual(serverCopy, other)
}

func ipsAreEqual(a, b []net.IP) (equal bool) {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if !a[i].Equal(b[i]) {
			return false
		}
	}

	return true
}

func (s *Server) Key() (key string) {
	var protocols []string
	if s.TCP {
		protocols = append(protocols, "tcp")
	}
	if s.UDP {
		protocols = append(protocols, "udp")
	}

	return fmt.Sprintf("%s-%s-%s", s.VPN, strings.Join(protocols, "-"), s.Hostname)
}
