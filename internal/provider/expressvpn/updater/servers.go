// package expressvpn contains code to obtain the server information
// for the ExpressVPN provider.
package expressvpn

import (
	"context"
	"fmt"
	"sort"

	"github.com/qdm12/gluetun/internal/constants/vpn"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/common"
)

func (u *Updater) GetServers(ctx context.Context, minServers int) (
	servers []models.Server, err error) {
	servers = hardcodedServers()

	hosts := make([]string, len(servers))
	for i := range servers {
		hosts[i] = servers[i].Hostname
	}

	hostToIPs, warnings, err := resolveHosts(ctx, u.presolver, hosts, minServers)
	for _, warning := range warnings {
		u.warner.Warn(warning)
	}
	if err != nil {
		return nil, err
	}

	i := 0
	for _, server := range servers {
		hostname := server.Hostname
		server.IPs = hostToIPs[hostname]
		if len(server.IPs) == 0 {
			continue
		}
		server.VPN = vpn.OpenVPN
		server.UDP = true // no TCP support
		servers[i] = server
		i++
	}
	servers = servers[:i]

	if len(servers) < minServers {
		return nil, fmt.Errorf("%w: %d and expected at least %d",
			common.ErrNotEnoughServers, len(servers), minServers)
	}

	sort.Sort(models.SortableServers(servers))

	return servers, nil
}
