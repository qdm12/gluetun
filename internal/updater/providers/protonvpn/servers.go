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
	servers []models.Server, warnings []string, err error) {
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

	ipToServer := make(ipToServer, count)
	for _, logicalServer := range data.LogicalServers {
		region := getStringValue(logicalServer.Region)
		city := getStringValue(logicalServer.City)
		name := logicalServer.Name
		for _, physicalServer := range logicalServer.Servers {
			if physicalServer.Status == 0 { // disabled so skip server
				warnings = append(warnings,
					"ignoring server "+physicalServer.Domain+" with status 0")
				continue
			}

			hostname := physicalServer.Domain
			entryIP := physicalServer.EntryIP

			// Note: for multi-hop use the server name or hostname
			// instead of the country
			countryCode := logicalServer.ExitCountry
			country, warning := codeToCountry(countryCode, countryCodes)
			if warning != "" {
				warnings = append(warnings, warning)
			}

			ipToServer.add(country, region, city, name, hostname, entryIP)
		}
	}

	if len(ipToServer) < minServers {
		return nil, warnings, fmt.Errorf("%w: %d and expected at least %d",
			ErrNotEnoughServers, len(ipToServer), minServers)
	}

	servers = ipToServer.toServersSlice()

	sortServers(servers)

	return servers, warnings, nil
}

func getStringValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}
