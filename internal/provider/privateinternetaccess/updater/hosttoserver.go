package updater

import (
	"net/netip"

	"github.com/qdm12/gluetun/internal/constants/vpn"
	"github.com/qdm12/gluetun/internal/models"
)

type nameToServer map[string][]models.Server

func (nts nameToServer) add(name, hostname, region string,
	serverType string, tcp, udp, portForward bool, ip netip.Addr) (change bool) {

	var server models.Server

	// Check for existing server for this name.
	servers := nts[name]
	for i, existingServer := range servers {
		if existingServer.ServerName != name || existingServer.VPN != serverType {
			continue
		}

		server = existingServer

		switch existingServer.VPN {
		case vpn.OpenVPN:
			// Update OpenVPN supported protocols and return
			if !existingServer.TCP {
				servers[i].TCP = tcp
				change = true
			}
			if !existingServer.UDP {
				servers[i].UDP = udp
				change = true
			}
			ipFound := false
			for _, existingIP := range server.IPs {
				if ip == existingIP {
					ipFound = true
					break
				}
			}
			if !ipFound {
				server.IPs = append(server.IPs, ip)
				change = true
			}
			break
		case vpn.Wireguard:
			// Update IPs and return
			ipFound := false
			for _, existingIP := range server.IPs {
				if ip == existingIP {
					ipFound = true
					break
				}
			}
			if !ipFound {
				server.IPs = append(server.IPs, ip)
				change = true
			}
			break
		}

		break
	}

	if server.ServerName == "" {
		change = true
		switch serverType {
		case vpn.OpenVPN:
			nts[name] = append(servers, models.Server{
				VPN:         vpn.OpenVPN,
				Region:      region,
				Hostname:    hostname,
				IPs:         []netip.Addr{ip},
				PortForward: portForward,
				ServerName:  name,
				TCP:         tcp,
				UDP:         udp,
			})
			break
		case vpn.Wireguard:
			nts[name] = append(servers, models.Server{
				VPN:         vpn.Wireguard,
				Region:      region,
				Hostname:    hostname,
				IPs:         []netip.Addr{ip},
				PortForward: portForward,
				ServerName:  name,
				UDP:         udp,
			})
			break
		}
	}

	return change
}

func (nts nameToServer) toServersSlice() (servers []models.Server) {
	servers = make([]models.Server, 0, len(nts))
	for _, hostServers := range nts {
		servers = append(servers, hostServers...)
	}
	return servers
}
