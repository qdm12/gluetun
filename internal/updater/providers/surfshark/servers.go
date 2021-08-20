// Package surfshark contains code to obtain the server information
// for the Surshark provider.
package surfshark

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/updater/resolver"
	"github.com/qdm12/gluetun/internal/updater/unzip"
)

var (
	ErrGetZip           = errors.New("cannot get OpenVPN ZIP file")
	ErrGetAPI           = errors.New("cannot fetch server information from API")
	ErrNotEnoughServers = errors.New("not enough servers found")
)

func GetServers(ctx context.Context, unzipper unzip.Unzipper,
	client *http.Client, presolver resolver.Parallel, minServers int) (
	servers []models.SurfsharkServer, warnings []string, err error) {
	hts := make(hostToServer)

	err = addServersFromAPI(ctx, client, hts)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %s", ErrGetAPI, err)
	}

	warnings, err = addOpenVPNServersFromZip(ctx, unzipper, hts)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %s", ErrGetZip, err)
	}

	getRemainingServers(hts)

	hosts := hts.toHostsSlice()
	hostToIPs, newWarnings, err := resolveHosts(ctx, presolver, hosts, minServers)
	warnings = append(warnings, newWarnings...)
	if err != nil {
		return nil, warnings, err
	}

	hts.adaptWithIPs(hostToIPs)

	servers = hts.toServersSlice()

	if len(servers) < minServers {
		return nil, warnings, fmt.Errorf("%w: %d and expected at least %d",
			ErrNotEnoughServers, len(servers), minServers)
	}

	sortServers(servers)

	return servers, warnings, nil
}
