// Package ivpn contains code to obtain the server information
// for the Surshark provider.
package ivpn

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/updater/resolver"
)

var (
	ErrFetchAPI         = errors.New("failed fetching API")
	ErrNotEnoughServers = errors.New("not enough servers found")
)

func GetServers(ctx context.Context, client *http.Client,
	presolver resolver.Parallel, minServers int) (
	servers []models.IvpnServer, warnings []string, err error) {
	data, err := fetchAPI(ctx, client)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %s", ErrFetchAPI, err)
	}

	hosts := make([]string, 0, len(data.Servers))

	for _, serverData := range data.Servers {
		host := serverData.Hostnames.OpenVPN

		if host == "" {
			continue // Wireguard
		}

		hosts = append(hosts, host)
	}

	if len(hosts) < minServers {
		return nil, nil, fmt.Errorf("%w: %d and expected at least %d",
			ErrNotEnoughServers, len(hosts), minServers)
	}

	hostToIPs, warnings, err := resolveHosts(ctx, presolver, hosts, minServers)
	if err != nil {
		return nil, warnings, err
	}

	servers = make([]models.IvpnServer, 0, len(hosts))
	for _, serverData := range data.Servers {
		host := serverData.Hostnames.OpenVPN
		if serverData.Hostnames.OpenVPN == "" {
			continue // Wireguard
		}

		server := models.IvpnServer{
			Country:  serverData.Country,
			City:     serverData.City,
			Hostname: serverData.Hostnames.OpenVPN,
			// TCP is not supported
			UDP: true,
			IPs: hostToIPs[host],
		}
		servers = append(servers, server)
	}

	sortServers(servers)

	return servers, warnings, nil
}
