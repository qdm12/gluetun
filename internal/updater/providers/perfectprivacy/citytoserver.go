package perfectprivacy

import (
	"net"

	"github.com/qdm12/gluetun/internal/models"
)

type cityToServer map[string]models.PerfectprivacyServer

func (cts cityToServer) add(city string, ips []net.IP) {
	server, ok := cts[city]
	if !ok {
		server.City = city
		server.IPs = ips
		server.TCP = true
		server.UDP = true
	} else {
		// Do not insert duplicate IP addresses
		existingIPs := make(map[string]struct{}, len(server.IPs))
		for _, ip := range server.IPs {
			existingIPs[ip.String()] = struct{}{}
		}

		for _, ip := range ips {
			ipString := ip.String()
			_, ok := existingIPs[ipString]
			if ok {
				continue
			}
			existingIPs[ipString] = struct{}{}
			server.IPs = append(server.IPs, ip)
		}
	}

	cts[city] = server
}

func (cts cityToServer) toServersSlice() (servers []models.PerfectprivacyServer) {
	servers = make([]models.PerfectprivacyServer, 0, len(cts))
	for _, server := range cts {
		servers = append(servers, server)
	}
	return servers
}
