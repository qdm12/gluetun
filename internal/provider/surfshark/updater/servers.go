// Package surfshark contains code to obtain the server information
// for the Surshark provider.
package surfshark

import (
	"context"
	"fmt"
	"sort"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/common"
)

func (u *Updater) GetServers(ctx context.Context, minServers int) (
	servers []models.Server, err error) {
	hts := make(hostToServer)

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
	hostToIPs, warnings, err := resolveHosts(ctx, u.presolver, hosts, minServers)
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
