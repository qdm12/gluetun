package updater

import (
	"net/netip"

	"github.com/qdm12/gluetun/internal/constants/vpn"
	"github.com/qdm12/gluetun/internal/models"
)

type ipToServer map[string]models.Server

type features struct {
	secureCore bool
	tor        bool
	p2p        bool
	stream     bool
}

func (its ipToServer) add(country, region, city, name, hostname string,
	free bool, entryIP netip.Addr, features features) {
	key := entryIP.String()

	server, ok := its[key]
	if ok {
		return
	}

	server.VPN = vpn.OpenVPN
	server.Country = country
	server.Region = region
	server.City = city
	server.ServerName = name
	server.Hostname = hostname
	server.Free = free
	server.SecureCore = features.secureCore
	server.Tor = features.tor
	server.PortForward = features.p2p
	server.Stream = features.stream
	server.UDP = true
	server.TCP = true
	server.IPs = []netip.Addr{entryIP}
	its[key] = server
}

func (its ipToServer) toServersSlice() (servers []models.Server) {
	servers = make([]models.Server, 0, len(its))
	for _, server := range its {
		servers = append(servers, server)
	}
	return servers
}
