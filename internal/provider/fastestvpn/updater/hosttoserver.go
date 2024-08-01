package updater

import (
	"net/netip"

	"github.com/qdm12/gluetun/internal/constants/vpn"
	"github.com/qdm12/gluetun/internal/models"
)

type hostToServerData map[string]serverData

type serverData struct {
	openvpn    bool
	wireguard  bool
	country    string
	city       string
	openvpnUDP bool
	openvpnTCP bool
	ips        []netip.Addr
}

func (hts hostToServerData) add(host, vpnType, country, city string, tcp, udp bool) {
	serverData, ok := hts[host]
	switch vpnType {
	case vpn.OpenVPN:
		serverData.openvpn = true
		serverData.openvpnTCP = serverData.openvpnTCP || tcp
		serverData.openvpnUDP = serverData.openvpnUDP || udp
	case vpn.Wireguard:
		serverData.wireguard = true
	default:
		panic("protocol not supported")
	}

	if !ok {
		serverData.country = country
		serverData.city = city
	} else if city != "" {
		// some servers are listed without the city although
		// they are also listed with the city described, so update
		// the city field.
		serverData.city = city
	}

	hts[host] = serverData
}

func (hts hostToServerData) toHostsSlice() (hosts []string) {
	hosts = make([]string, 0, len(hts))
	for host := range hts {
		hosts = append(hosts, host)
	}
	return hosts
}

func (hts hostToServerData) adaptWithIPs(hostToIPs map[string][]netip.Addr) {
	for host, serverData := range hts {
		ips := hostToIPs[host]
		if len(ips) == 0 {
			delete(hts, host)
			continue
		}
		serverData.ips = ips
		hts[host] = serverData
	}
}

func (hts hostToServerData) toServersSlice() (servers []models.Server) {
	servers = make([]models.Server, 0, 2*len(hts)) //nolint:gomnd
	for hostname, serverData := range hts {
		server := models.Server{
			VPN:      vpn.Both,
			Hostname: hostname,
			Country:  serverData.country,
			City:     serverData.city,
			IPs:      serverData.ips,
			TCP:      serverData.openvpnTCP,
			UDP:      serverData.openvpnUDP,
			WgPubKey: "658QxufMbjOTmB61Z7f+c7Rjg7oqWLnepTalqBERjF0=",
		}
		servers = append(servers, server)
	}
	return servers
}
