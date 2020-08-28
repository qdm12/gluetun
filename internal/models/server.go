package models

import (
	"fmt"
	"net"
	"strings"
)

func stringifyIPs(ips []net.IP) string {
	ipStrings := make([]string, len(ips))
	for i := range ips {
		ipStrings[i] = fmt.Sprintf("{%s}", strings.ReplaceAll(ips[i].String(), ".", ", "))
	}
	return strings.Join(ipStrings, ", ")
}

type PIAServer struct {
	IPs    []net.IP `json:"ips"`
	Region string   `json:"region"`
}

func (p *PIAServer) String() string {
	return fmt.Sprintf("{Region: %q, IPs: []net.IP{%s}}", p.Region, stringifyIPs(p.IPs))
}

type MullvadServer struct {
	IPs     []net.IP `json:"ips"`
	Country string   `json:"country"`
	City    string   `json:"city"`
	ISP     string   `json:"isp"`
	Owned   bool     `json:"owned"`
}

type WindscribeServer struct {
	Region string   `json:"region"`
	IPs    []net.IP `json:"ips"`
}

type SurfsharkServer struct {
	Region string   `json:"region"`
	IPs    []net.IP `json:"ips"`
}

type CyberghostServer struct {
	Region string   `json:"region"`
	Group  string   `json:"group"`
	IPs    []net.IP `json:"ips"`
}

type VyprvpnServer struct {
	Region string   `json:"region"`
	IPs    []net.IP `json:"ips"`
}

type NordvpnServer struct { //nolint:maligned
	Region string `json:"region"`
	Number uint16 `json:"number"`
	IP     net.IP `json:"ip"`
	TCP    bool   `json:"tcp"`
	UDP    bool   `json:"udp"`
}

type PurevpnServer struct {
	Region  string   `json:"region"`
	Country string   `json:"country"`
	City    string   `json:"city"`
	IPs     []net.IP `json:"ips"`
}
