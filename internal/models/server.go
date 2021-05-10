package models

import (
	"encoding/hex"
	"fmt"
	"net"
	"strings"
)

type CyberghostServer struct {
	Region   string   `json:"region"`
	Group    string   `json:"group"`
	Hostname string   `json:"hostname"`
	IPs      []net.IP `json:"ips"`
}

func (s *CyberghostServer) String() string {
	return fmt.Sprintf("{Region: %q, Group: %q, Hostname: %q, IPs: %s}",
		s.Region, s.Group, s.Hostname, goStringifyIPs(s.IPs))
}

type FastestvpnServer struct {
	Hostname string   `json:"hostname"`
	TCP      bool     `json:"tcp"`
	UDP      bool     `json:"udp"`
	Country  string   `json:"country"`
	IPs      []net.IP `json:"ips"`
}

func (s *FastestvpnServer) String() string {
	return fmt.Sprintf("{Country: %q, Hostname: %q, UDP: %t, TCP: %t, IPs: %s}",
		s.Country, s.Hostname, s.UDP, s.TCP, goStringifyIPs(s.IPs))
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

func (s *HideMyAssServer) String() string {
	return fmt.Sprintf("{Country: %q, Region: %q, City: %q, Hostname: %q, TCP: %t, UDP: %t, IPs: %s}",
		s.Country, s.Region, s.City, s.Hostname, s.TCP, s.UDP, goStringifyIPs(s.IPs))
}

type MullvadServer struct {
	IPs      []net.IP `json:"ips"`
	IPsV6    []net.IP `json:"ipsv6"`
	Country  string   `json:"country"`
	City     string   `json:"city"`
	Hostname string   `json:"hostname"`
	ISP      string   `json:"isp"`
	Owned    bool     `json:"owned"`
}

func (s *MullvadServer) String() string {
	return fmt.Sprintf("{Country: %q, City: %q, Hostname: %q, ISP: %q, Owned: %t, IPs: %s, IPsV6: %s}",
		s.Country, s.City, s.Hostname, s.ISP, s.Owned, goStringifyIPs(s.IPs), goStringifyIPs(s.IPsV6))
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

func (s *NordvpnServer) String() string {
	return fmt.Sprintf("{Region: %q, Hostname: %q, Name: %q, Number: %d, TCP: %t, UDP: %t, IP: %s}",
		s.Region, s.Hostname, s.Name, s.Number, s.TCP, s.UDP, goStringifyIP(s.IP))
}

type PrivadoServer struct {
	Country  string `json:"country"`
	Region   string `json:"region"`
	City     string `json:"city"`
	Hostname string `json:"hostname"`
	IP       net.IP `json:"ip"`
}

func (s *PrivadoServer) String() string {
	return fmt.Sprintf("{Country: %q, Region: %q, City: %q, Hostname: %q, IP: %s}",
		s.Country, s.Region, s.City, s.Hostname, goStringifyIP(s.IP))
}

type PIAServer struct {
	Region      string `json:"region"`
	Hostname    string `json:"hostname"`
	ServerName  string `json:"server_name"`
	TCP         bool   `json:"tcp"`
	UDP         bool   `json:"udp"`
	PortForward bool   `json:"port_forward"`
	IP          net.IP `json:"ip"`
}

func (p *PIAServer) String() string {
	return fmt.Sprintf("{Region: %q, Hostname: %q, ServerName: %q, TCP: %t, UDP: %t, PortForward: %t, IP: %s}",
		p.Region, p.Hostname, p.ServerName, p.TCP, p.UDP, p.PortForward, goStringifyIP(p.IP))
}

type PrivatevpnServer struct {
	Country  string   `json:"country"`
	City     string   `json:"city"`
	Hostname string   `json:"hostname"`
	IPs      []net.IP `json:"ip"`
}

func (s *PrivatevpnServer) String() string {
	return fmt.Sprintf("{Country: %q, City: %q, Hostname: %q, IPs: %s}",
		s.Country, s.City, s.Hostname, goStringifyIPs(s.IPs))
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

func (s *ProtonvpnServer) String() string {
	return fmt.Sprintf("{Country: %q, Region: %q, City: %q, Name: %q, Hostname: %q, EntryIP: %s, ExitIP: %s}",
		s.Country, s.Region, s.City, s.Name, s.Hostname, goStringifyIP(s.EntryIP), goStringifyIP(s.ExitIP))
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

func (s *PurevpnServer) String() string {
	return fmt.Sprintf("{Country: %q, Region: %q, City: %q, Hostname: %q, TCP: %t, UDP: %t, IPs: %s}",
		s.Country, s.Region, s.City, s.Hostname, s.TCP, s.UDP, goStringifyIPs(s.IPs))
}

type SurfsharkServer struct {
	Region string   `json:"region"`
	IPs    []net.IP `json:"ips"`
}

func (s *SurfsharkServer) String() string {
	return fmt.Sprintf("{Region: %q, IPs: %s}", s.Region, goStringifyIPs(s.IPs))
}

type TorguardServer struct {
	Country  string `json:"country"`
	City     string `json:"city"`
	Hostname string `json:"hostname"`
	IP       net.IP `json:"ip"`
}

func (s *TorguardServer) String() string {
	return fmt.Sprintf("{Country: %q, City: %q, Hostname: %q, IP: %s}",
		s.Country, s.City, s.Hostname, goStringifyIP(s.IP))
}

type VyprvpnServer struct {
	Region string   `json:"region"`
	IPs    []net.IP `json:"ips"`
}

func (s *VyprvpnServer) String() string {
	return fmt.Sprintf("{Region: %q, IPs: %s}", s.Region, goStringifyIPs(s.IPs))
}

type WindscribeServer struct {
	Region   string `json:"region"`
	City     string `json:"city"`
	Hostname string `json:"hostname"`
	IP       net.IP `json:"ip"`
}

func (s *WindscribeServer) String() string {
	return fmt.Sprintf("{Region: %q, City: %q, Hostname: %q, IP: %s}",
		s.Region, s.City, s.Hostname, goStringifyIP(s.IP))
}

func goStringifyIP(ip net.IP) string {
	s := fmt.Sprintf("%#v", ip)
	s = strings.TrimSuffix(strings.TrimPrefix(s, "net.IP{"), "}")
	fields := strings.Split(s, ", ")
	isIPv4 := ip.To4() != nil
	if isIPv4 {
		fields = fields[len(fields)-4:]
	}

	// Count leading zeros
	leadingZeros := 0
	for i := range fields {
		if fields[i] == "0x0" {
			leadingZeros++
		} else {
			break
		}
	}

	// Remove leading zeros
	fields = fields[leadingZeros:]

	for i := range fields {
		// IPv4 is better understood in integer notation, whereas IPv6 is written in hex notation
		if isIPv4 {
			hexString := strings.Replace(fields[i], "0x", "", 1)
			if len(hexString) == 1 {
				hexString = "0" + hexString
			}
			b, _ := hex.DecodeString(hexString)
			fields[i] = fmt.Sprintf("%d", b[0])
		}
	}

	return fmt.Sprintf("net.IP{%s}", strings.Join(fields, ", "))
}

func goStringifyIPs(ips []net.IP) string {
	ipStrings := make([]string, len(ips))
	for i := range ips {
		ipStrings[i] = goStringifyIP(ips[i])
		ipStrings[i] = strings.TrimPrefix(ipStrings[i], "net.IP")
	}
	return "[]net.IP{" + strings.Join(ipStrings, ", ") + "}"
}
