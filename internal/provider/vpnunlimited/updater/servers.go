// Package vpnunlimited contains code to obtain the server information
// for the VPNUnlimited provider.
package vpnunlimited

import (
	"context"
	"fmt"
	"sort"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/common"
)

func (u *Updater) GetServers(ctx context.Context, minServers int) (
	servers []models.Server, err error) {
	// Hardcoded data from a user provided ZIP file since it's behind a login wall
	hts, warnings := getHostToServer()
	for _, warning := range warnings {
		u.warner.Warn(warning)
	}

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
