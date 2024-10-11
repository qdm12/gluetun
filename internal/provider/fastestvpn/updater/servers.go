package updater

import (
	"context"
	"fmt"
	"sort"

	"github.com/qdm12/gluetun/internal/constants/vpn"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/common"
)

func (u *Updater) FetchServers(ctx context.Context, minServers int) (
	servers []models.Server, err error,
) {
	protocols := []string{"ikev2", "tcp", "udp"}
	hts := make(hostToServerData)

	for _, protocol := range protocols {
		apiServers, err := fetchAPIServers(ctx, u.client, protocol)
		if err != nil {
			return nil, fmt.Errorf("fetching %s servers from API: %w", protocol, err)
		}
		for _, apiServer := range apiServers {
			// all hostnames from the protocols TCP, UDP and IKEV2 support Wireguard
			// per https://github.com/qdm12/gluetun-wiki/issues/76#issuecomment-2125420536
			const wgTCP, wgUDP = false, false // ignored
			hts.add(apiServer.hostname, vpn.Wireguard, apiServer.country, apiServer.city, wgTCP, wgUDP)

			tcp := protocol == "tcp"
			udp := protocol == "udp"
			if !tcp && !udp { // not an OpenVPN protocol, for example ikev2
				continue
			}
			hts.add(apiServer.hostname, vpn.OpenVPN, apiServer.country, apiServer.city, tcp, udp)
		}
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

	hts.adaptWithIPs(hostToIPs)

	servers = hts.toServersSlice()

	if len(servers) < minServers {
		return nil, fmt.Errorf("%w: %d and expected at least %d",
			common.ErrNotEnoughServers, len(servers), minServers)
	}

	sort.Sort(models.SortableServers(servers))

	return servers, nil
}
