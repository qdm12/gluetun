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
	hts := make(hostToServers)

	err = addServersFromAPI(ctx, u.client, hts)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch server information from API: %w", err)
	}

	warnings, err := addOpenVPNServersFromZip(ctx, u.unzipper, hts)
	for _, warning := range warnings {
		u.warner.Warn(warning)
	}
	if err != nil {
		return nil, fmt.Errorf("cannot get OpenVPN ZIP file: %w", err)
	}

	getRemainingServers(hts)

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

	if len(hts) < minServers {
		return nil, fmt.Errorf("%w: %d and expected at least %d",
			common.ErrNotEnoughServers, len(hts), minServers)
	}

	servers = hts.toServersSlice()

	sort.Sort(models.SortableServers(servers))

	return servers, nil
}
