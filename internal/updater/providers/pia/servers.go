// Package pia contains code to obtain the server information
// for the Private Internet Access provider.
package pia

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/qdm12/gluetun/internal/models"
)

var ErrNotEnoughServers = errors.New("not enough servers found")

func GetServers(ctx context.Context, client *http.Client, minServers int) (
	servers []models.PIAServer, err error) {
	data, err := fetchAPI(ctx, client)
	if err != nil {
		return nil, err
	}

	nts := make(nameToServer)

	for _, region := range data.Regions {
		for _, server := range region.Servers.UDP {
			const tcp, udp = false, true
			nts.add(server.CN, region.DNS, region.Name, tcp, udp, region.PortForward, server.IP)
		}

		for _, server := range region.Servers.TCP {
			const tcp, udp = true, false
			nts.add(server.CN, region.DNS, region.Name, tcp, udp, region.PortForward, server.IP)
		}
	}

	servers = nts.toServersSlice()

	if len(servers) < minServers {
		return nil, fmt.Errorf("%w: %d and expected at least %d",
			ErrNotEnoughServers, len(servers), minServers)
	}

	sortServers(servers)

	return servers, nil
}
