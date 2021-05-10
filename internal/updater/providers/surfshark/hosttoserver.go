package surfshark

import (
	"net"
	"strings"

	"github.com/qdm12/gluetun/internal/models"
)

type hostToServer map[string]models.SurfsharkServer

func (hts hostToServer) add(host, region string, tcp, udp bool) {
	server, ok := hts[host]
	if !ok {
		server.Hostname = host
		server.Region = region
	}
	if tcp {
		server.TCP = tcp
	}
	if udp {
		server.UDP = udp
	}
	hts[host] = server
}

func (hts hostToServer) toHostsSlice() (hosts []string) {
	hosts = make([]string, 0, len(hts))
	for host := range hts {
		hosts = append(hosts, host)
	}
	return hosts
}

func (hts hostToServer) toSubdomainsSlice() (subdomains []string) {
	subdomains = make([]string, 0, len(hts))
	const suffix = ".prod.surfshark.com"
	for host := range hts {
		subdomain := strings.TrimSuffix(host, suffix)
		subdomains = append(subdomains, subdomain)
	}
	return subdomains
}

func (hts hostToServer) adaptWithIPs(hostToIPs map[string][]net.IP) {
	for host, IPs := range hostToIPs {
		server := hts[host]
		server.IPs = IPs
		hts[host] = server
	}
	for host, server := range hts {
		if len(server.IPs) == 0 {
			delete(hts, host)
		}
	}
}

func (hts hostToServer) toServersSlice() (servers []models.SurfsharkServer) {
	servers = make([]models.SurfsharkServer, 0, len(hts))
	for _, server := range hts {
		servers = append(servers, server)
	}
	return servers
}
