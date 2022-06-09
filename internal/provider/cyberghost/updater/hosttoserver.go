package updater

import (
	"net"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/constants/vpn"
	"github.com/qdm12/gluetun/internal/models"
)

type hostToServer map[string]models.Server

func getPossibleServers() (possibleServers hostToServer) {
	groupIDToProtocol := getGroupIDToProtocol()

	cyberghostCountryCodes := getSubdomainToRegion()
	allCountryCodes := constants.CountryCodes()
	possibleCountryCodes := mergeCountryCodes(cyberghostCountryCodes, allCountryCodes)

	n := len(groupIDToProtocol) * len(possibleCountryCodes)

	possibleServers = make(hostToServer, n) // key is the host

	for groupID, protocol := range groupIDToProtocol {
		for countryCode, country := range possibleCountryCodes {
			const domain = "cg-dialup.net"
			possibleHost := groupID + "-" + countryCode + "." + domain
			possibleServer := models.Server{
				VPN:      vpn.OpenVPN,
				Hostname: possibleHost,
				Country:  country,
				TCP:      protocol == constants.TCP,
				UDP:      protocol == constants.UDP,
			}
			possibleServers[possibleHost] = possibleServer
		}
	}

	return possibleServers
}

func (hts hostToServer) hostsSlice() (hosts []string) {
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

func (hts hostToServer) toSlice() (servers []models.Server) {
	servers = make([]models.Server, 0, len(hts))
	for _, server := range hts {
		servers = append(servers, server)
	}
	return servers
}
