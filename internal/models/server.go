package models

import (
	"encoding/hex"
	"fmt"
	"net"
	"strings"
)

type PIAServer struct {
	Region      string           `json:"region"`
	PortForward bool             `json:"port_forward"`
	OpenvpnUDP  PIAServerOpenvpn `json:"openvpn_udp"`
	OpenvpnTCP  PIAServerOpenvpn `json:"openvpn_tcp"`
}

type PIAServerOpenvpn struct {
	IPs []net.IP `json:"ips"`
	CN  string   `json:"cn"`
}

func (p *PIAServerOpenvpn) String() string {
	return fmt.Sprintf("models.PIAServerOpenvpn{CN: %q, IPs: %s}", p.CN, goStringifyIPs(p.IPs))
}

func (p *PIAServer) String() string {
	return fmt.Sprintf("{Region: %q, PortForward: %t, OpenvpnUDP: %s, OpenvpnTCP: %s}",
		p.Region, p.PortForward, p.OpenvpnUDP.String(), p.OpenvpnTCP.String())
}

type MullvadServer struct {
	IPs     []net.IP `json:"ips"`
	IPsV6   []net.IP `json:"ipsv6"`
	Country string   `json:"country"`
	City    string   `json:"city"`
	ISP     string   `json:"isp"`
	Owned   bool     `json:"owned"`
}

func (s *MullvadServer) String() string {
	return fmt.Sprintf("{Country: %q, City: %q, ISP: %q, Owned: %t, IPs: %s, IPsV6: %s}",
		s.Country, s.City, s.ISP, s.Owned, goStringifyIPs(s.IPs), goStringifyIPs(s.IPsV6))
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

type SurfsharkServer struct {
	Region string   `json:"region"`
	IPs    []net.IP `json:"ips"`
}

func (s *SurfsharkServer) String() string {
	return fmt.Sprintf("{Region: %q, IPs: %s}", s.Region, goStringifyIPs(s.IPs))
}

type CyberghostServer struct {
	Region string   `json:"region"`
	Group  string   `json:"group"`
	IPs    []net.IP `json:"ips"`
}

func (s *CyberghostServer) String() string {
	return fmt.Sprintf("{Region: %q, Group: %q, IPs: %s}", s.Region, s.Group, goStringifyIPs(s.IPs))
}

type VyprvpnServer struct {
	Region string   `json:"region"`
	IPs    []net.IP `json:"ips"`
}

func (s *VyprvpnServer) String() string {
	return fmt.Sprintf("{Region: %q, IPs: %s}", s.Region, goStringifyIPs(s.IPs))
}

type NordvpnServer struct { //nolint:maligned
	Region string `json:"region"`
	Number uint16 `json:"number"`
	IP     net.IP `json:"ip"`
	TCP    bool   `json:"tcp"`
	UDP    bool   `json:"udp"`
}

func (s *NordvpnServer) String() string {
	return fmt.Sprintf("{Region: %q, Number: %d, TCP: %t, UDP: %t, IP: %s}",
		s.Region, s.Number, s.TCP, s.UDP, goStringifyIP(s.IP))
}

type PurevpnServer struct {
	Region  string   `json:"region"`
	Country string   `json:"country"`
	City    string   `json:"city"`
	IPs     []net.IP `json:"ips"`
}

func (s *PurevpnServer) String() string {
	return fmt.Sprintf("{Region: %q, Country: %q, City: %q, IPs: %s}",
		s.Region, s.Country, s.City, goStringifyIPs(s.IPs))
}

type PrivadoServer struct {
	IP     net.IP `json:"ip"`
	City   string `json:"city"`
	Number uint16 `json:"number"`
}

func (s *PrivadoServer) String() string {
	return fmt.Sprintf("{City: %q, Number: %d, IP: %s}",
		s.City, s.Number, goStringifyIP(s.IP))
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
