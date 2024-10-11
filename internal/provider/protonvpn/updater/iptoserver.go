package updater

import (
	"net/netip"

	"github.com/qdm12/gluetun/internal/constants/vpn"
	"github.com/qdm12/gluetun/internal/models"
)

type ipToServers map[string][2]models.Server // first server is OpenVPN, second is Wireguard.

type features struct {
	secureCore bool
	tor        bool
	p2p        bool
	stream     bool
}

func (its ipToServers) add(country, region, city, name, hostname, wgPubKey string,
	free bool, entryIP netip.Addr, features features,
) {
	key := entryIP.String()

	servers, ok := its[key]
	if ok {
		return
	}

	baseServer := models.Server{
		Country:     country,
		Region:      region,
		City:        city,
		ServerName:  name,
		Hostname:    hostname,
		Free:        free,
		SecureCore:  features.secureCore,
		Tor:         features.tor,
		PortForward: features.p2p,
		Stream:      features.stream,
		IPs:         []netip.Addr{entryIP},
	}
	openvpnServer := baseServer
	openvpnServer.VPN = vpn.OpenVPN
	openvpnServer.UDP = true
	openvpnServer.TCP = true
	servers[0] = openvpnServer
	wireguardServer := baseServer
	wireguardServer.VPN = vpn.Wireguard
	wireguardServer.WgPubKey = wgPubKey
	servers[1] = wireguardServer
	its[key] = servers
}

func (its ipToServers) toServersSlice() (serversSlice []models.Server) {
	const vpnProtocols = 2
	serversSlice = make([]models.Server, 0, vpnProtocols*len(its))
	for _, servers := range its {
		serversSlice = append(serversSlice, servers[0], servers[1])
	}
	return serversSlice
}

func (its ipToServers) numberOfServers() int {
	const serversPerIP = 2
	return len(its) * serversPerIP
}
