package updater

import (
	"net/netip"

	"github.com/qdm12/gluetun/internal/constants/vpn"
	"github.com/qdm12/gluetun/internal/models"
)

type hostToServer map[string]models.Server

func (hts hostToServer) add(host, vpnType, region, country, city,
	retroLoc, wgPubKey string, openvpnTCP, openvpnUDP bool) {
	server, ok := hts[host]
	if !ok {
		server := models.Server{
			VPN:      vpnType,
			Region:   region,
			Country:  country,
			City:     city,
			RetroLoc: retroLoc,
			Hostname: host,
			WgPubKey: wgPubKey,
			TCP:      openvpnTCP,
			UDP:      openvpnUDP,
		}
		hts[host] = server
		return
	}

	server.SetVPN(vpnType)
	if vpnType == vpn.OpenVPN {
		server.TCP = server.TCP || openvpnTCP
		server.UDP = server.UDP || openvpnUDP
	} else if wgPubKey != "" {
		server.WgPubKey = wgPubKey
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

func (hts hostToServer) adaptWithIPs(hostToIPs map[string][]netip.Addr) {
	for host, server := range hts {
		ips := hostToIPs[host]
		if len(ips) == 0 {
			delete(hts, host)
			continue
		}
		server.IPs = ips
		hts[host] = server
	}
}

func (hts hostToServer) toServersSlice() (servers []models.Server) {
	servers = make([]models.Server, 0, len(hts))
	for _, server := range hts {
		servers = append(servers, server)
	}
	return servers
}
