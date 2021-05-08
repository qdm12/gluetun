// Package protonvpn contains code to obtain the server information
// for the ProtonVPN provider.
package protonvpn

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
)

var ErrNotEnoughServers = errors.New("not enough servers found")

func GetServers(ctx context.Context, client *http.Client, minServers int) (
	servers []models.ProtonvpnServer, warnings []string, err error) {
	data, err := fetchAPI(ctx, client)
	if err != nil {
		return nil, nil, err
	}

	countryCodes := constants.CountryCodes()

	var count int
	for _, logicalServer := range data.LogicalServers {
		count += len(logicalServer.Servers)
	}

	if count < minServers {
		return nil, warnings, fmt.Errorf("%w: %d and expected at least %d",
			ErrNotEnoughServers, count, minServers)
	}

	servers = make([]models.ProtonvpnServer, 0, count)
	for _, logicalServer := range data.LogicalServers {
		for _, physicalServer := range logicalServer.Servers {
			server, warning, err := makeServer(
				physicalServer, logicalServer, countryCodes)

			if warning != "" {
				warnings = append(warnings, warning)
			}

			if err != nil {
				warnings = append(warnings, err.Error())
				continue
			}

			servers = append(servers, server)
		}
	}

	if len(servers) < minServers {
		return nil, warnings, fmt.Errorf("%w: %d and expected at least %d",
			ErrNotEnoughServers, len(servers), minServers)
	}

	sortServers(servers)

	return servers, warnings, nil
}

var errServerStatusZero = errors.New("ignoring server with status 0")

func makeServer(physical physicalServer, logical logicalServer,
	countryCodes map[string]string) (server models.ProtonvpnServer,
	warning string, err error) {
	if physical.Status == 0 {
		return server, "", fmt.Errorf("%w: %s",
			errServerStatusZero, physical.Domain)
	}

	countryCode := logical.ExitCountry
	country, warning := codeToCountry(countryCode, countryCodes)

	server = models.ProtonvpnServer{
		// Note: for multi-hop use the server name or hostname
		// instead of the country
		Country:  country,
		Region:   getStringValue(logical.Region),
		City:     getStringValue(logical.City),
		Name:     logical.Name,
		Hostname: physical.Domain,
		EntryIP:  physical.EntryIP,
		ExitIP:   physical.ExitIP,
	}

	return server, warning, nil
}

func getStringValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}
