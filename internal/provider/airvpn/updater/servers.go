package updater

import (
	"context"
	"fmt"
	"net"
	"sort"

	"github.com/qdm12/gluetun/internal/constants/vpn"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/common"
)

func (u *Updater) FetchServers(ctx context.Context, minServers int) (
	servers []models.Server, err error) {
	data, err := fetchAPI(ctx, u.client)
	if err != nil {
		return nil, fmt.Errorf("fetching API: %w", err)
	}

	// every API server model has:
	// - Wireguard server using IPv4In1
	// - Wiregard server using IPv6In1
	// - OpenVPN TCP+UDP server using IPv4In1
	// - OpenVPN TCP+UDP server using IPv4In3
	// - OpenVPN TCP+UDP server using IPv6In1
	// - OpenVPN TCP+UDP server using IPv6In3
	// - OpenVPN UDP server using IPv4In2
	// - OpenVPN UDP server using IPv4In4
	// - OpenVPN UDP server using IPv6In2
	// - OpenVPN UDP server using IPv6In4
	const numberOfServersPerAPIServer = 10
	projectedNumberOfServers := numberOfServersPerAPIServer * len(data.Servers)

	if projectedNumberOfServers < minServers {
		return nil, fmt.Errorf("%w: %d and expected at least %d",
			common.ErrNotEnoughServers, projectedNumberOfServers, minServers)
	}

	servers = make([]models.Server, 0, projectedNumberOfServers)
	for _, apiServer := range data.Servers {
		if apiServer.Health != "ok" {
			continue
		}

		baseServer := models.Server{
			ServerName: apiServer.PublicName,
			Country:    apiServer.CountryName,
			City:       apiServer.Location,
			Region:     apiServer.Continent,
		}

		baseWireguardServer := baseServer
		baseWireguardServer.VPN = vpn.Wireguard
		baseWireguardServer.WgPubKey = "PyLCXAQT8KkM4T+dUsOQfn+Ub3pGxfGlxkIApuig+hk="

		ipv4WireguadServer := baseWireguardServer
		ipv4WireguadServer.IPs = []net.IP{apiServer.IPv4In1}
		ipv4WireguadServer.Hostname = apiServer.CountryCode + ".vpn.airdns.org"
		servers = append(servers, ipv4WireguadServer)

		ipv6WireguadServer := baseWireguardServer
		ipv6WireguadServer.IPs = []net.IP{apiServer.IPv6In1}
		ipv6WireguadServer.Hostname = apiServer.CountryCode + ".ipv6.vpn.airdns.org"
		servers = append(servers, ipv6WireguadServer)

		baseOpenVPNServer := baseServer
		baseOpenVPNServer.VPN = vpn.OpenVPN
		baseOpenVPNServer.UDP = true

		ipv4In2OpenVPNServer := baseOpenVPNServer
		ipv4In2OpenVPNServer.IPs = []net.IP{apiServer.IPv4In2}
		ipv4In2OpenVPNServer.Hostname = apiServer.CountryCode + "2.vpn.airdns.org"
		servers = append(servers, ipv4In2OpenVPNServer)

		ipv6In2OpenVPNServer := baseOpenVPNServer
		ipv6In2OpenVPNServer.IPs = []net.IP{apiServer.IPv6In2}
		ipv6In2OpenVPNServer.Hostname = apiServer.CountryCode + "2.ipv6.vpn.airdns.org"
		servers = append(servers, ipv6In2OpenVPNServer)

		// TODO uncomment when AirVPN clarifies why these don't work for OpenVPN UDP on port 1194
		// ipv4In4OpenVPNServer := baseOpenVPNServer
		// ipv4In4OpenVPNServer.IPs = []net.IP{apiServer.IPv4In4}
		// ipv4In4OpenVPNServer.Hostname = apiServer.CountryCode + "4.vpn.airdns.org"
		// servers = append(servers, ipv4In4OpenVPNServer)

		// ipv6In4OpenVPNServer := baseOpenVPNServer
		// ipv6In4OpenVPNServer.IPs = []net.IP{apiServer.IPv6In4}
		// ipv6In4OpenVPNServer.Hostname = apiServer.CountryCode + "4.ipv6.vpn.airdns.org"
		// servers = append(servers, ipv6In4OpenVPNServer)

		baseOpenVPNServer.TCP = true // both TCP+UDP servers below

		ipv4In1OpenVPNServer := baseOpenVPNServer
		ipv4In1OpenVPNServer.IPs = []net.IP{apiServer.IPv4In1}
		ipv4In1OpenVPNServer.Hostname = apiServer.CountryCode + ".vpn.airdns.org"
		servers = append(servers, ipv4In1OpenVPNServer)

		ipv6In1OpenVPNServer := baseOpenVPNServer
		ipv6In1OpenVPNServer.IPs = []net.IP{apiServer.IPv6In1}
		ipv6In1OpenVPNServer.Hostname = apiServer.CountryCode + ".ipv6.vpn.airdns.org"
		servers = append(servers, ipv6In1OpenVPNServer)

		// TODO this should support TCP but it only supports UDP.
		baseOpenVPNServer.TCP = false
		ipv4In3OpenVPNServer := baseOpenVPNServer
		ipv4In3OpenVPNServer.IPs = []net.IP{apiServer.IPv4In3}
		ipv4In3OpenVPNServer.Hostname = apiServer.CountryCode + "3.vpn.airdns.org"
		servers = append(servers, ipv4In3OpenVPNServer)

		ipv6In3OpenVPNServer := baseOpenVPNServer
		ipv6In3OpenVPNServer.IPs = []net.IP{apiServer.IPv6In3}
		ipv6In3OpenVPNServer.Hostname = apiServer.CountryCode + "3.ipv6.vpn.airdns.org"
		servers = append(servers, ipv6In3OpenVPNServer)
	}

	sort.Sort(models.SortableServers(servers))

	return servers, nil
}
