// Package mullvad contains code to obtain the server information
// for the Mullvad provider.
package mullvad

import (
	"context"
	"fmt"
	"sort"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/common"
)

func (u *Updater) GetServers(ctx context.Context, minServers int) (
	servers []models.Server, err error) {
	data, err := fetchAPI(ctx, u.client)
	if err != nil {
		return nil, err
	}

	hts := make(hostToServer)
	for _, serverData := range data {
		if err := hts.add(serverData); err != nil {
			return nil, err
		}
	}

	if len(hts) < minServers {
		return nil, fmt.Errorf("%w: %d and expected at least %d",
			common.ErrNotEnoughServers, len(hts), minServers)
	}

	servers = hts.toServersSlice()

	sort.Sort(models.SortableServers(servers))

	return servers, nil
}
