// Package mullvad contains code to obtain the server information
// for the Mullvad provider.
package mullvad

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/qdm12/gluetun/internal/models"
)

var ErrNotEnoughServers = errors.New("not enough servers found")

func GetServers(ctx context.Context, client *http.Client, minServers int) (
	servers []models.Server, err error) {
	data, err := fetchAPI(ctx, client)
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
			ErrNotEnoughServers, len(hts), minServers)
	}

	servers = hts.toServersSlice()

	sortServers(servers)

	return servers, nil
}
