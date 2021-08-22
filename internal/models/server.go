package models

import (
	"net"
)

type CyberghostServer struct {
	Region   string   `json:"region"`
	Group    string   `json:"group"`
	Hostname string   `json:"hostname"`
	IPs      []net.IP `json:"ips"`
}

type FastestvpnServer struct {
	Hostname string   `json:"hostname"`
	TCP      bool     `json:"tcp"`
	UDP      bool     `json:"udp"`
	Country  string   `json:"country"`
	IPs      []net.IP `json:"ips"`
}

type HideMyAssServer struct {
	Country  string   `json:"country"`
	Region   string   `json:"region"`
	City     string   `json:"city"`
	Hostname string   `json:"hostname"`
	TCP      bool     `json:"tcp"`
	UDP      bool     `json:"udp"`
	IPs      []net.IP `json:"ips"`
}

type IpvanishServer struct {
	Country  string   `json:"country"`
	City     string   `json:"city"`
	Hostname string   `json:"hostname"`
	TCP      bool     `json:"tcp"`
	UDP      bool     `json:"udp"`
	IPs      []net.IP `json:"ips"`
}

type IvpnServer struct {
	Country  string   `json:"country"`
	City     string   `json:"city"`
	Hostname string   `json:"hostname"`
	TCP      bool     `json:"tcp"`
	UDP      bool     `json:"udp"`
	IPs      []net.IP `json:"ips"`
}

type MullvadServer struct {
	VPN      string   `json:"vpn"`
	IPs      []net.IP `json:"ips"`
	IPsV6    []net.IP `json:"ipsv6"`
	Country  string   `json:"country"`
	City     string   `json:"city"`
	Hostname string   `json:"hostname"`
	ISP      string   `json:"isp"`
	Owned    bool     `json:"owned"`
	WgPubKey string   `json:"wgpubkey,omitempty"`
}

type NordvpnServer struct { //nolint:maligned
	Region   string `json:"region"`
	Hostname string `json:"hostname"`
	Name     string `json:"name"`
	Number   uint16 `json:"number"`
	IP       net.IP `json:"ip"`
	TCP      bool   `json:"tcp"`
	UDP      bool   `json:"udp"`
}

type PrivadoServer struct {
	Country  string `json:"country"`
	Region   string `json:"region"`
	City     string `json:"city"`
	Hostname string `json:"hostname"`
	IP       net.IP `json:"ip"`
}

type PIAServer struct {
	Region      string   `json:"region"`
	Hostname    string   `json:"hostname"`
	ServerName  string   `json:"server_name"`
	TCP         bool     `json:"tcp"`
	UDP         bool     `json:"udp"`
	PortForward bool     `json:"port_forward"`
	IPs         []net.IP `json:"ips"`
}

type PrivatevpnServer struct {
	Country  string   `json:"country"`
	City     string   `json:"city"`
	Hostname string   `json:"hostname"`
	IPs      []net.IP `json:"ip"`
}

type ProtonvpnServer struct {
	Country  string `json:"country"`
	Region   string `json:"region"`
	City     string `json:"city"`
	Name     string `json:"name"`
	Hostname string `json:"hostname"`
	EntryIP  net.IP `json:"entry_ip"`
	ExitIP   net.IP `json:"exit_ip"` // TODO verify it matches with public IP once connected
}

type PurevpnServer struct {
	Country  string   `json:"country"`
	Region   string   `json:"region"`
	City     string   `json:"city"`
	Hostname string   `json:"hostname"`
	TCP      bool     `json:"tcp"`
	UDP      bool     `json:"udp"`
	IPs      []net.IP `json:"ips"`
}

type SurfsharkServer struct {
	Region   string   `json:"region"`
	Hostname string   `json:"hostname"`
	TCP      bool     `json:"tcp"`
	UDP      bool     `json:"udp"`
	IPs      []net.IP `json:"ips"`
}

type TorguardServer struct {
	Country  string   `json:"country"`
	City     string   `json:"city"`
	Hostname string   `json:"hostname"`
	TCP      bool     `json:"tcp"`
	UDP      bool     `json:"udp"`
	IPs      []net.IP `json:"ips"`
}

type VPNUnlimitedServer struct {
	Country  string   `json:"country"`
	City     string   `json:"city"`
	Hostname string   `json:"hostname"`
	Free     bool     `json:"free"`
	Stream   bool     `json:"stream"`
	TCP      bool     `json:"tcp"`
	UDP      bool     `json:"udp"`
	IPs      []net.IP `json:"ips"`
}

type VyprvpnServer struct {
	Region   string   `json:"region"`
	Hostname string   `json:"hostname"`
	TCP      bool     `json:"tcp"`
	UDP      bool     `json:"udp"` // only support for UDP
	IPs      []net.IP `json:"ips"`
}

type WindscribeServer struct {
	VPN      string   `json:"vpn"`
	Region   string   `json:"region"`
	City     string   `json:"city"`
	Hostname string   `json:"hostname"`
	OvpnX509 string   `json:"x509,omitempty"`
	WgPubKey string   `json:"wgpubkey,omitempty"`
	IPs      []net.IP `json:"ips"`
}
