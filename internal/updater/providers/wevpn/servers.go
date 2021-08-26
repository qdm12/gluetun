// package wevpn contains code to obtain the server information
// for the WeVPN provider.
package wevpn

import (
	"context"
	"errors"
	"fmt"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/updater/resolver"
)

var (
	ErrGetZip           = errors.New("cannot get OpenVPN ZIP file")
	ErrGetAPI           = errors.New("cannot fetch server information from API")
	ErrNotEnoughServers = errors.New("not enough servers found")
)

func GetServers(ctx context.Context, presolver resolver.Parallel, minServers int) (
	servers []models.WevpnServer, warnings []string, err error) {
	cities := getAvailableCities()
	servers = make([]models.WevpnServer, 0, len(cities))
	hostnames := make([]string, len(cities))
	hostnameToCity := make(map[string]string, len(cities))

	for i, city := range cities {
		hostname := getHostnameFromCity(city)
		hostnames[i] = hostname
		hostnameToCity[hostname] = city
	}

	hostnameToIPs, newWarnings, err := resolveHosts(ctx, presolver, hostnames, minServers)
	warnings = append(warnings, newWarnings...)
	if err != nil {
		return nil, warnings, err
	}

	if len(hostnameToIPs) < minServers {
		return nil, warnings, fmt.Errorf("%w: %d and expected at least %d",
			ErrNotEnoughServers, len(servers), minServers)
	}

	for hostname, ips := range hostnameToIPs {
		city := hostnameToCity[hostname]
		server := models.WevpnServer{
			City:     city,
			Hostname: hostname,
			IPs:      ips,
		}
		servers = append(servers, server)
	}

	sortServers(servers)

	return servers, warnings, nil
}
