package protonvpn

import (
	"net"

	"github.com/qdm12/gluetun/internal/models"
)

type ipToServer map[string]models.ProtonvpnServer

func (its ipToServer) add(country, region, city, name, hostname string,
	entryIP net.IP) {
	key := entryIP.String()

	server, ok := its[key]
	if !ok {
		server.Country = country
		server.Region = region
		server.City = city
		server.Name = name
		server.Hostname = hostname
		server.IPs = []net.IP{entryIP}
	} else {
		server.IPs = append(server.IPs, entryIP)
	}

	its[key] = server
}

func (its ipToServer) toServersSlice() (servers []models.ProtonvpnServer) {
	servers = make([]models.ProtonvpnServer, 0, len(its))
	for _, server := range its {
		servers = append(servers, server)
	}
	return servers
}
