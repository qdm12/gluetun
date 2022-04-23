// Package hidemyass contains code to obtain the server information
// for the HideMyAss provider.
package hidemyass

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/qdm12/gluetun/internal/constants/vpn"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/updater/resolver"
)

var ErrNotEnoughServers = errors.New("not enough servers found")

func GetServers(ctx context.Context, client *http.Client,
	presolver resolver.Parallel, minServers int) (
	servers []models.Server, warnings []string, err error) {
	tcpHostToURL, udpHostToURL, err := getAllHostToURL(ctx, client)
	if err != nil {
		return nil, nil, err
	}

	hosts := getUniqueHosts(tcpHostToURL, udpHostToURL)

	if len(hosts) < minServers {
		return nil, nil, fmt.Errorf("%w: %d and expected at least %d",
			ErrNotEnoughServers, len(hosts), minServers)
	}

	hostToIPs, warnings, err := resolveHosts(ctx, presolver, hosts, minServers)
	if err != nil {
		return nil, warnings, err
	}

	servers = make([]models.Server, 0, len(hostToIPs))
	for host, IPs := range hostToIPs {
		tcpURL, tcp := tcpHostToURL[host]
		udpURL, udp := udpHostToURL[host]

		// These two are only used to extract the country, region and city.
		var url, protocol string
		if tcp {
			url = tcpURL
			protocol = "TCP"
		} else if udp {
			url = udpURL
			protocol = "UDP"
		}
		country, region, city := parseOpenvpnURL(url, protocol)

		server := models.Server{
			VPN:      vpn.OpenVPN,
			Country:  country,
			Region:   region,
			City:     city,
			Hostname: host,
			IPs:      IPs,
			TCP:      tcp,
			UDP:      udp,
		}
		servers = append(servers, server)
	}

	sortServers(servers)

	return servers, warnings, nil
}
