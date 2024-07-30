package updater

import (
	"context"
	"fmt"
	"sort"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/common"
)

func (u *Updater) FetchServers(ctx context.Context, minServers int) (
	servers []models.Server, err error) {
	protocols := []string{"tcp", "udp"}
	hts := make(hostToServer)

	for _, protocol := range protocols {
		apiServers, err := fetchAPIServers(ctx, u.client, protocol)
		if err != nil {
			return nil, fmt.Errorf("fetching %s servers from API: %w", protocol, err)
		}
		for _, apiServer := range apiServers {
			tcp := protocol == "tcp"
			udp := protocol == "udp"
			hts.add(apiServer.hostname, apiServer.country, apiServer.city, tcp, udp)
		}
	}

	if len(hts) < minServers {
		return nil, fmt.Errorf("%w: %d and expected at least %d",
			common.ErrNotEnoughServers, len(hts), minServers)
	}

	hosts := hts.toHostsSlice()
	resolveSettings := parallelResolverSettings(hosts)
	hostToIPs, warnings, err := u.parallelResolver.Resolve(ctx, resolveSettings)
	for _, warning := range warnings {
		u.warner.Warn(warning)
	}
	if err != nil {
		return nil, err
	}

	if len(hostToIPs) < minServers {
		return nil, fmt.Errorf("%w: %d and expected at least %d",
			common.ErrNotEnoughServers, len(servers), minServers)
	}

	hts.adaptWithIPs(hostToIPs)

	servers = hts.toServersSlice()

	sort.Sort(models.SortableServers(servers))

	return servers, nil
}
