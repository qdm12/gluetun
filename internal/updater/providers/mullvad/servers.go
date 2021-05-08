// Package mullvad contains code to obtain the server information
// for the Mullvad provider.
package mullvad

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/qdm12/gluetun/internal/models"
)

var ErrNotEnoughServers = errors.New("not enough servers found")

func GetServers(ctx context.Context, client *http.Client, minServers int) (
	servers []models.MullvadServer, err error) {
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

	servers = groupByProperties(servers)

	sortServers(servers)

	return servers, nil
}

// TODO group by hostname so remove this.
func groupByProperties(serversByHost []models.MullvadServer) (serversByProps []models.MullvadServer) {
	propsToServer := make(map[string]models.MullvadServer, len(serversByHost))
	for _, server := range serversByHost {
		key := server.Country + server.City + server.ISP + strconv.FormatBool(server.Owned)
		serverByProps, ok := propsToServer[key]
		if !ok {
			serverByProps.Country = server.Country
			serverByProps.City = server.City
			serverByProps.ISP = server.ISP
			serverByProps.Owned = server.Owned
		}
		serverByProps.IPs = append(serverByProps.IPs, server.IPs...)
		serverByProps.IPsV6 = append(serverByProps.IPsV6, server.IPsV6...)
		propsToServer[key] = serverByProps
	}

	serversByProps = make([]models.MullvadServer, 0, len(propsToServer))
	for _, serverByProp := range propsToServer {
		serversByProps = append(serversByProps, serverByProp)
	}

	return serversByProps
}
