// package expressvpn contains code to obtain the server information
// for the ExpressVPN provider.
package expressvpn

import (
	"context"
	"errors"
	"fmt"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/updater/resolver"
	"github.com/qdm12/gluetun/internal/updater/unzip"
)

var ErrNotEnoughServers = errors.New("not enough servers found")

func GetServers(ctx context.Context, unzipper unzip.Unzipper,
	presolver resolver.Parallel, minServers int) (
	servers []models.ExpressvpnServer, warnings []string, err error) {
	servers = hardcodedServers()

	hosts := make([]string, len(servers))
	for i := range servers {
		hosts[i] = servers[i].Hostname
	}

	hostToIPs, newWarnings, err := resolveHosts(ctx, presolver, hosts, minServers)
	warnings = append(warnings, newWarnings...)
	if err != nil {
		return nil, warnings, err
	}

	i := 0
	for _, server := range servers {
		hostname := server.Hostname
		server.IPs = hostToIPs[hostname]
		if len(server.IPs) == 0 {
			continue
		}
		server.UDP = true // no TCP support
		servers[i] = server
		i++
	}
	servers = servers[:i]

	if len(servers) < minServers {
		return nil, warnings, fmt.Errorf("%w: %d and expected at least %d",
			ErrNotEnoughServers, len(servers), minServers)
	}

	sortServers(servers)

	return servers, warnings, nil
}
