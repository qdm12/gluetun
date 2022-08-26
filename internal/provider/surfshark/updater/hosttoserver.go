package updater

import (
	"net"

	"github.com/qdm12/gluetun/internal/constants/vpn"
	"github.com/qdm12/gluetun/internal/models"
)

type hostToServers map[string][]models.Server

func (hts hostToServers) addOpenVPN(host, region, country, city,
	retroLoc string, tcp, udp bool) {
	// Check for existing server for this host and OpenVPN.
	servers := hts[host]
	for i, existingServer := range servers {
		if existingServer.Hostname != host ||
			existingServer.VPN != vpn.OpenVPN {
			continue
		}

		// Update OpenVPN supported protocols and return
		if !existingServer.TCP {
			servers[i].TCP = tcp
		}
		if !existingServer.UDP {
			servers[i].UDP = udp
		}
		return
	}

	server := models.Server{
		VPN:      vpn.OpenVPN,
		Region:   region,
		Country:  country,
		City:     city,
		RetroLoc: retroLoc,
		Hostname: host,
		TCP:      tcp,
		UDP:      udp,
	}
	hts[host] = append(servers, server)
}

func (hts hostToServers) addWireguard(host, region, country, city, retroLoc,
	wgPubKey string) {
	// Check for existing server for this host and Wireguard.
	servers := hts[host]
	for _, existingServer := range servers {
		if existingServer.Hostname == host &&
			existingServer.VPN == vpn.Wireguard {
			// No update necessary for Wireguard
			return
		}
	}

	server := models.Server{
		VPN:      vpn.Wireguard,
		Region:   region,
		Country:  country,
		City:     city,
		RetroLoc: retroLoc,
		Hostname: host,
		WgPubKey: wgPubKey,
	}
	hts[host] = append(servers, server)
}

func (hts hostToServers) toHostsSlice() (hosts []string) {
	const vpnServerTypes = 2 // OpenVPN + Wireguard
	hosts = make([]string, 0, vpnServerTypes*len(hts))
	for host := range hts {
		hosts = append(hosts, host)
	}
	return hosts
}

func (hts hostToServers) adaptWithIPs(hostToIPs map[string][]net.IP) {
	for host, IPs := range hostToIPs {
		servers := hts[host]
		for i := range servers {
			servers[i].IPs = IPs
		}
		hts[host] = servers
	}
	for host, servers := range hts {
		if len(servers[0].IPs) == 0 {
			delete(hts, host)
		}
	}
}

func (hts hostToServers) toServersSlice() (servers []models.Server) {
	const vpnServerTypes = 2 // OpenVPN + Wireguard
	servers = make([]models.Server, 0, vpnServerTypes*len(hts))
	for _, serversForHost := range hts {
		servers = append(servers, serversForHost...)
	}
	return servers
}
