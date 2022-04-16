package models

import (
	"net"
)

type CyberghostServer struct {
	Country  string   `json:"country,omitempty"`
	Hostname string   `json:"hostname,omitempty"`
	TCP      bool     `json:"tcp,omitempty"`
	UDP      bool     `json:"udp,omitempty"`
	IPs      []net.IP `json:"ips,omitempty"`
}

type ExpressvpnServer struct {
	Country  string   `json:"country,omitempty"`
	City     string   `json:"city,omitempty"`
	Hostname string   `json:"hostname,omitempty"`
	TCP      bool     `json:"tcp,omitempty"`
	UDP      bool     `json:"udp,omitempty"`
	IPs      []net.IP `json:"ips,omitempty"`
}

type FastestvpnServer struct {
	Hostname string   `json:"hostname,omitempty"`
	TCP      bool     `json:"tcp,omitempty"`
	UDP      bool     `json:"udp,omitempty"`
	Country  string   `json:"country,omitempty"`
	IPs      []net.IP `json:"ips,omitempty"`
}

type HideMyAssServer struct {
	Country  string   `json:"country,omitempty"`
	Region   string   `json:"region,omitempty"`
	City     string   `json:"city,omitempty"`
	Hostname string   `json:"hostname,omitempty"`
	TCP      bool     `json:"tcp,omitempty"`
	UDP      bool     `json:"udp,omitempty"`
	IPs      []net.IP `json:"ips,omitempty"`
}

type IpvanishServer struct {
	Country  string   `json:"country,omitempty"`
	City     string   `json:"city,omitempty"`
	Hostname string   `json:"hostname,omitempty"`
	TCP      bool     `json:"tcp,omitempty"`
	UDP      bool     `json:"udp,omitempty"`
	IPs      []net.IP `json:"ips,omitempty"`
}

type IvpnServer struct {
	VPN      string   `json:"vpn,omitempty"`
	Country  string   `json:"country,omitempty"`
	City     string   `json:"city,omitempty"`
	ISP      string   `json:"isp,omitempty"`
	Hostname string   `json:"hostname,omitempty"`
	WgPubKey string   `json:"wgpubkey,omitempty"`
	TCP      bool     `json:"tcp,omitempty"`
	UDP      bool     `json:"udp,omitempty"`
	IPs      []net.IP `json:"ips,omitempty"`
}

type MullvadServer struct {
	VPN      string   `json:"vpn,omitempty"`
	IPs      []net.IP `json:"ips,omitempty"`
	IPsV6    []net.IP `json:"ipsv6,omitempty"`
	Country  string   `json:"country,omitempty"`
	City     string   `json:"city,omitempty"`
	Hostname string   `json:"hostname,omitempty"`
	ISP      string   `json:"isp,omitempty"`
	Owned    bool     `json:"owned,omitempty"`
	WgPubKey string   `json:"wgpubkey,omitempty"`
}

type NordvpnServer struct { //nolint:maligned
	Region   string `json:"region,omitempty"`
	Hostname string `json:"hostname,omitempty"`
	Number   uint16 `json:"number,omitempty"`
	IP       net.IP `json:"ip,omitempty"`
	TCP      bool   `json:"tcp,omitempty"`
	UDP      bool   `json:"udp,omitempty"`
}

type PerfectprivacyServer struct {
	City string   `json:"city,omitempty"` // primary key
	IPs  []net.IP `json:"ips,omitempty"`
	TCP  bool     `json:"tcp,omitempty"`
	UDP  bool     `json:"udp,omitempty"`
}

type PrivadoServer struct {
	Country  string `json:"country,omitempty"`
	Region   string `json:"region,omitempty"`
	City     string `json:"city,omitempty"`
	Hostname string `json:"hostname,omitempty"`
	IP       net.IP `json:"ip,omitempty"`
}

type PIAServer struct {
	Region      string   `json:"region,omitempty"`
	Hostname    string   `json:"hostname,omitempty"`
	ServerName  string   `json:"server_name,omitempty"`
	TCP         bool     `json:"tcp,omitempty"`
	UDP         bool     `json:"udp,omitempty"`
	PortForward bool     `json:"port_forward,omitempty"`
	IPs         []net.IP `json:"ips,omitempty"`
}

type PrivatevpnServer struct {
	Country  string   `json:"country,omitempty"`
	City     string   `json:"city,omitempty"`
	Hostname string   `json:"hostname,omitempty"`
	IPs      []net.IP `json:"ip,omitempty"`
}

type ProtonvpnServer struct {
	Country  string `json:"country,omitempty"`
	Region   string `json:"region,omitempty"`
	City     string `json:"city,omitempty"`
	Name     string `json:"server_name,omitempty"`
	Hostname string `json:"hostname,omitempty"`
	EntryIP  net.IP `json:"entry_ip,omitempty"`
}

type PurevpnServer struct {
	Country  string   `json:"country,omitempty"`
	Region   string   `json:"region,omitempty"`
	City     string   `json:"city,omitempty"`
	Hostname string   `json:"hostname,omitempty"`
	TCP      bool     `json:"tcp,omitempty"`
	UDP      bool     `json:"udp,omitempty"`
	IPs      []net.IP `json:"ips,omitempty"`
}

type SurfsharkServer struct {
	Region   string   `json:"region,omitempty"`
	Country  string   `json:"country,omitempty"` // Country is also used for multi-hop
	City     string   `json:"city,omitempty"`
	RetroLoc string   `json:"retroloc,omitempty"` // TODO remove in v4
	Hostname string   `json:"hostname,omitempty"`
	MultiHop bool     `json:"multihop,omitempty"`
	TCP      bool     `json:"tcp,omitempty"`
	UDP      bool     `json:"udp,omitempty"`
	IPs      []net.IP `json:"ips,omitempty"`
}

type TorguardServer struct {
	Country  string   `json:"country,omitempty"`
	City     string   `json:"city,omitempty"`
	Hostname string   `json:"hostname,omitempty"`
	TCP      bool     `json:"tcp,omitempty"`
	UDP      bool     `json:"udp,omitempty"`
	IPs      []net.IP `json:"ips,omitempty"`
}

type VPNUnlimitedServer struct {
	Country  string   `json:"country,omitempty"`
	City     string   `json:"city,omitempty"`
	Hostname string   `json:"hostname,omitempty"`
	Free     bool     `json:"free,omitempty"`
	Stream   bool     `json:"stream,omitempty"`
	TCP      bool     `json:"tcp,omitempty"`
	UDP      bool     `json:"udp,omitempty"`
	IPs      []net.IP `json:"ips,omitempty"`
}

type VyprvpnServer struct {
	Region   string   `json:"region,omitempty"`
	Hostname string   `json:"hostname,omitempty"`
	TCP      bool     `json:"tcp,omitempty"`
	UDP      bool     `json:"udp,omitempty"` // only support for UDP
	IPs      []net.IP `json:"ips,omitempty"`
}

type WevpnServer struct {
	City     string   `json:"city,omitempty"`
	Hostname string   `json:"hostname,omitempty"`
	TCP      bool     `json:"tcp,omitempty"`
	UDP      bool     `json:"udp,omitempty"`
	IPs      []net.IP `json:"ips,omitempty"`
}

type WindscribeServer struct {
	VPN      string   `json:"vpn,omitempty"`
	Region   string   `json:"region,omitempty"`
	City     string   `json:"city,omitempty"`
	Hostname string   `json:"hostname,omitempty"`
	OvpnX509 string   `json:"x509,omitempty"`
	WgPubKey string   `json:"wgpubkey,omitempty"`
	IPs      []net.IP `json:"ips,omitempty"`
}
