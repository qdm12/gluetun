package purevpn

import (
	"net"

	"github.com/qdm12/gluetun/internal/models"
)

type locationToServer map[string]models.PurevpnServer

func locationKey(country, region, city string) string {
	return country + region + city
}

func (lts locationToServer) add(country, region, city string, ips []net.IP) {
	key := locationKey(country, region, city)
	server, ok := lts[key]
	if !ok {
		server.Country = country
		server.Region = region
		server.City = city
	}
	server.IPs = append(server.IPs, ips...)
	lts[key] = server
}

func (lts locationToServer) toServersSlice() (servers []models.PurevpnServer) {
	servers = make([]models.PurevpnServer, 0, len(lts))
	for _, server := range lts {
		servers = append(servers, server)
	}
	return servers
}
