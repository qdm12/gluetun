// package wevpn contains code to obtain the server information
// for the WeVPN provider.
package wevpn

import (
	"context"
	"errors"
	"fmt"

	"github.com/qdm12/gluetun/internal/constants/vpn"
	"github.com/qdm12/gluetun/internal/models"
)

var (
	ErrGetZip           = errors.New("cannot get OpenVPN ZIP file")
	ErrGetAPI           = errors.New("cannot fetch server information from API")
	ErrNotEnoughServers = errors.New("not enough servers found")
)

func (u *Updater) GetServers(ctx context.Context, minServers int) (
	servers []models.Server, err error) {
	cities := getAvailableCities()
	servers = make([]models.Server, 0, len(cities))
	hostnames := make([]string, len(cities))
	hostnameToCity := make(map[string]string, len(cities))

	for i, city := range cities {
		hostname := getHostnameFromCity(city)
		hostnames[i] = hostname
		hostnameToCity[hostname] = city
	}

	hostnameToIPs, warnings, err := resolveHosts(ctx, u.presolver, hostnames, minServers)
	for _, warning := range warnings {
		u.warner.Warn(warning)
	}
	if err != nil {
		return nil, err
	}

	if len(hostnameToIPs) < minServers {
		return nil, fmt.Errorf("%w: %d and expected at least %d",
			ErrNotEnoughServers, len(servers), minServers)
	}

	for hostname, ips := range hostnameToIPs {
		city := hostnameToCity[hostname]
		server := models.Server{
			VPN:      vpn.OpenVPN,
			City:     city,
			Hostname: hostname,
			UDP:      true,
			IPs:      ips,
		}
		servers = append(servers, server)
	}

	sortServers(servers)

	return servers, nil
}
