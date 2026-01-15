package updater

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/common"
)

func (u *Updater) FetchServers(ctx context.Context, minServers int) (
	servers []models.Server, err error,
) {
	possibleServers := getPossibleServers()

	possibleHosts := possibleServers.hostsSlice()
	resolveSettings := parallelResolverSettings(possibleHosts)
	hostToIPs, warnings, err := u.parallelResolver.Resolve(ctx, resolveSettings)
	for _, warning := range warnings {
		if strings.HasSuffix(warning, "no such host") {
			continue // ignore no such host warnings
		}
		u.warner.Warn(warning)
	}
	if err != nil {
		return nil, err
	}

	if len(hostToIPs) < minServers {
		return nil, fmt.Errorf("%w: %d and expected at least %d",
			common.ErrNotEnoughServers, len(servers), minServers)
	}

	possibleServers.adaptWithIPs(hostToIPs)

	servers = possibleServers.toSlice()

	sort.Sort(models.SortableServers(servers))

	return servers, nil
}
