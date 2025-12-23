package updater

import (
	"context"
	"fmt"
	"sort"

	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/common"
)

func (u *Updater) FetchServers(ctx context.Context, minServers int) (
	servers []models.Server, err error,
) {
	// Since SlickVPN website listing VPN servers https://www.slickvpn.com/locations/
	// went to become a pile of trash, we now use the servers data from our servers.json
	// to check which servers can be resolved. The previous code was dynamically parsing
	// their website table of servers and they now list only 11 servers on their website.
	hardcodedServersData := u.storage.HardcodedServers()
	slickVPNData, ok := hardcodedServersData.ProviderToServers[providers.SlickVPN]
	if !ok {
		return nil, fmt.Errorf("no hardcoded servers for provider %s", providers.SlickVPN)
	}
	hardcodedServers := make([]models.Server, len(slickVPNData.Servers))
	copy(hardcodedServers, slickVPNData.Servers)

	hosts := make([]string, len(hardcodedServers))
	for i := range hardcodedServers {
		hosts[i] = hardcodedServers[i].Hostname
	}

	resolveSettings := parallelResolverSettings(hosts)
	hostToIPs, warnings, err := u.parallelResolver.Resolve(ctx, resolveSettings)
	for _, warning := range warnings {
		u.warner.Warn(warning)
	}
	if err != nil {
		return nil, fmt.Errorf("resolving hosts: %w", err)
	}

	if len(hostToIPs) < minServers {
		return nil, fmt.Errorf("%w: %d and expected at least %d",
			common.ErrNotEnoughServers, len(hosts), minServers)
	}

	servers = make([]models.Server, 0, len(hostToIPs))
	for _, server := range hardcodedServers {
		IPs, ok := hostToIPs[server.Hostname]
		if !ok || len(IPs) == 0 {
			continue
		}
		server.IPs = IPs
		servers = append(servers, server)
	}

	sort.Sort(models.SortableServers(servers))

	return servers, nil
}
