// Package cyberghost contains code to obtain the server information
// for the Cyberghost provider.
package cyberghost

import (
	"context"
	"sort"

	"github.com/qdm12/gluetun/internal/models"
)

func (u *Updater) GetServers(ctx context.Context, minServers int) (
	servers []models.Server, err error) {
	possibleServers := getPossibleServers()

	possibleHosts := possibleServers.hostsSlice()
	hostToIPs, err := resolveHosts(ctx, u.presolver, possibleHosts, minServers)
	if err != nil {
		return nil, err
	}

	possibleServers.adaptWithIPs(hostToIPs)

	servers = possibleServers.toSlice()

	sort.Sort(models.SortableServers(servers))

	return servers, nil
}
