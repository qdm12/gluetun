// Package nordvpn contains code to obtain the server information
// for the NordVPN provider.
package nordvpn

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/qdm12/gluetun/internal/models"
)

var (
	ErrParseIP          = errors.New("cannot parse IP address")
	ErrNotIPv4          = errors.New("IP address is not IPv4")
	ErrNotEnoughServers = errors.New("not enough servers found")
)

func GetServers(ctx context.Context, client *http.Client, minServers int) (
	servers []models.NordvpnServer, warnings []string, err error) {
	data, err := fetchAPI(ctx, client)
	if err != nil {
		return nil, nil, err
	}

	servers = make([]models.NordvpnServer, 0, len(data))

	for _, jsonServer := range data {
		if !jsonServer.Features.TCP && !jsonServer.Features.UDP {
			warning := "server does not support TCP and UDP for openvpn: " + jsonServer.Name
			warnings = append(warnings, warning)
			continue
		}

		ip, err := parseIPv4(jsonServer.IPAddress)
		if err != nil {
			return nil, nil, fmt.Errorf("%w for server %s", err, jsonServer.Name)
		}

		number, err := parseServerName(jsonServer.Name)
		if err != nil {
			return nil, nil, err
		}

		server := models.NordvpnServer{
			Region:   jsonServer.Country,
			Hostname: jsonServer.Domain,
			Name:     jsonServer.Name,
			Number:   number,
			IP:       ip,
			TCP:      jsonServer.Features.TCP,
			UDP:      jsonServer.Features.UDP,
		}
		servers = append(servers, server)
	}

	if len(servers) < minServers {
		return nil, warnings, fmt.Errorf("%w: %d and expected at least %d",
			ErrNotEnoughServers, len(servers), minServers)
	}

	sortServers(servers)

	return servers, warnings, nil
}
