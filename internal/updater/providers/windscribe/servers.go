// Package windscribe contains code to obtain the server information
// for the Windscribe provider.
package windscribe

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/qdm12/gluetun/internal/models"
)

var ErrNotEnoughServers = errors.New("not enough servers found")

func GetServers(ctx context.Context, client *http.Client, minServers int) (
	servers []models.WindscribeServer, err error) {
	data, err := fetchAPI(ctx, client)
	if err != nil {
		return nil, err
	}

	for _, regionData := range data.Data {
		region := regionData.Region
		for _, group := range regionData.Groups {
			city := group.City
			x5090Name := group.OvpnX509
			for _, node := range group.Nodes {
				const maxIPsPerNode = 3
				ips := make([]net.IP, 0, maxIPsPerNode)
				if node.IP != nil {
					ips = append(ips, node.IP)
				}
				if node.IP2 != nil {
					ips = append(ips, node.IP2)
				}
				if node.IP3 != nil {
					ips = append(ips, node.IP3)
				}
				server := models.WindscribeServer{
					Region:   region,
					City:     city,
					Hostname: node.Hostname,
					OvpnX509: x5090Name,
					IPs:      ips,
				}
				servers = append(servers, server)
			}
		}
	}

	if len(servers) < minServers {
		return nil, fmt.Errorf("%w: %d and expected at least %d",
			ErrNotEnoughServers, len(servers), minServers)
	}

	sortServers(servers)

	return servers, nil
}
