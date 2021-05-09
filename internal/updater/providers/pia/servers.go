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

	for _, region := range data.Regions {
		// Deduplicate servers with the same common name
		commonNameToProtocols := dedupByProtocol(region)

		// newServers can support only UDP or both TCP and UDP
		newServers := dataToServers(region.Servers.UDP, region.Name,
			region.DNS, region.PortForward, commonNameToProtocols)
		servers = append(servers, newServers...)

		// tcpServers only support TCP as mixed servers were found above.
		tcpServers := dataToServers(region.Servers.TCP, region.Name,
			region.DNS, region.PortForward, commonNameToProtocols)
		servers = append(servers, tcpServers...)
	}

	if len(servers) < minServers {
		return nil, fmt.Errorf("%w: %d and expected at least %d",
			ErrNotEnoughServers, len(servers), minServers)
	}

	sortServers(servers)

	return servers, nil
}

type protocols struct {
	tcp bool
	udp bool
}

// Deduplicate servers with the same common name for different protocols.
func dedupByProtocol(region regionData) (commonNameToProtocols map[string]protocols) {
	commonNameToProtocols = make(map[string]protocols)
	for _, udpServer := range region.Servers.UDP {
		protocols := commonNameToProtocols[udpServer.CN]
		protocols.udp = true
		commonNameToProtocols[udpServer.CN] = protocols
	}
	for _, tcpServer := range region.Servers.TCP {
		protocols := commonNameToProtocols[tcpServer.CN]
		protocols.tcp = true
		commonNameToProtocols[tcpServer.CN] = protocols
	}
	return commonNameToProtocols
}

func dataToServers(data []serverData, region, hostname string,
	portForward bool, commonNameToProtocols map[string]protocols) (
	servers []models.PIAServer) {
	servers = make([]models.PIAServer, 0, len(data))
	for _, serverData := range data {
		proto, ok := commonNameToProtocols[serverData.CN]
		if !ok {
			continue // server already added
		}
		delete(commonNameToProtocols, serverData.CN)
		server := models.PIAServer{
			Region:      region,
			Hostname:    hostname,
			ServerName:  serverData.CN,
			TCP:         proto.tcp,
			UDP:         proto.udp,
			PortForward: portForward,
			IP:          serverData.IP,
		}
		servers = append(servers, server)
	}
	return servers
}
