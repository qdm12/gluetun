// Package cyberghost contains code to obtain the server information
// for the Cyberghost provider.
package cyberghost

import (
	"context"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/updater/resolver"
)

func GetServers(ctx context.Context, presolver resolver.Parallel,
	minServers int) (servers []models.CyberghostServer, err error) {
	possibleServers := getPossibleServers()

	possibleHosts := possibleServers.hostsSlice()
	hostToIPs, err := resolveHosts(ctx, presolver, possibleHosts, minServers)
	if err != nil {
		return nil, err
	}

	possibleServers.adaptWithIPs(hostToIPs)

	servers = possibleServers.toSlice()

	sortServers(servers)
	return servers, nil
}
