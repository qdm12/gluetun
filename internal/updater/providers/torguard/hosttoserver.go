package torguard

import (
	"net"

	"github.com/qdm12/gluetun/internal/models"
)

type hostToServer map[string]models.TorguardServer

func (hts hostToServer) add(host, country, city string, tcp, udp bool, ip net.IP) {
	server, ok := hts[host]
	if !ok {
		server.Hostname = host
		server.Country = country
		server.City = city
		server.IPs = append(server.IPs, ip) // used if DNS resolution fails downstream
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

func (hts hostToServer) toServersSlice() (servers []models.TorguardServer) {
	servers = make([]models.TorguardServer, 0, len(hts))
	for _, server := range hts {
		servers = append(servers, server)
	}
	return servers
}
