package privado

import (
	"net"

	"github.com/qdm12/gluetun/internal/models"
)

type hostToServer map[string]models.PrivadoServer

func (hts hostToServer) add(host string) {
	server, ok := hts[host]
	if ok {
		return
	}
	server.Hostname = host
	hts[host] = server
}

func (hts hostToServer) toHostsSlice() (hosts []string) {
	hosts = make([]string, 0, len(hts))
	for host := range hts {
		hosts = append(hosts, host)
	}
	return hosts
}

func (hts hostToServer) adaptWithIPs(hostToIPs map[string][]net.IP) (
	warnings []string) {
	for host, IPs := range hostToIPs {
		if len(IPs) > 1 {
			warning := "more than one IP address found for host " + host
			warnings = append(warnings, warning)
		}
		server := hts[host]
		server.IP = IPs[0]
		hts[host] = server
	}
	for host, server := range hts {
		if server.IP == nil {
			delete(hts, host)
		}
	}
	return warnings
}

func (hts hostToServer) toServersSlice() (servers []models.PrivadoServer) {
	servers = make([]models.PrivadoServer, 0, len(hts))
	for _, server := range hts {
		servers = append(servers, server)
	}
	return servers
}
