package updater

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sort"
	"strings"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
)

func (u *updater) updateProtonvpn(ctx context.Context) (err error) {
	servers, warnings, err := findProtonvpnServers(ctx, u.client)
	if u.options.CLI {
		for _, warning := range warnings {
			u.logger.Warn("Protonvpn: %s", warning)
		}
	}
	if err != nil {
		return fmt.Errorf("cannot update Protonvpn servers: %w", err)
	}
	if u.options.Stdout {
		u.println(stringifyProtonvpnServers(servers))
	}
	u.servers.Protonvpn.Timestamp = u.timeNow().Unix()
	u.servers.Protonvpn.Servers = servers
	return nil
}

func findProtonvpnServers(ctx context.Context, client *http.Client) (
	servers []models.ProtonvpnServer, warnings []string, err error) {
	const url = "https://api.protonmail.ch/vpn/logicals"

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("%w: %s for %s", ErrHTTPStatusCodeNotOK, response.Status, url)
	}

	decoder := json.NewDecoder(response.Body)
	var data struct {
		LogicalServers []struct {
			Name        string
			ExitCountry string
			Region      *string
			City        *string
			Servers     []struct {
				EntryIP net.IP
				ExitIP  net.IP
				Domain  string
				Status  uint8
			}
		}
	}
	if err := decoder.Decode(&data); err != nil {
		return nil, nil, err
	}

	if err := response.Body.Close(); err != nil {
		return nil, nil, err
	}

	countryCodesMapping := constants.CountryCodes()
	for _, logicalServer := range data.LogicalServers {
		for _, physicalServer := range logicalServer.Servers {
			if physicalServer.Status == 0 {
				warnings = append(warnings, "ignoring server "+physicalServer.Domain+" as its status is 0")
				continue
			}

			countryCode := strings.ToLower(logicalServer.ExitCountry)
			country, ok := countryCodesMapping[countryCode]
			if !ok {
				warnings = append(warnings, "country not found for country code "+countryCode)
				country = logicalServer.ExitCountry
			}

			server := models.ProtonvpnServer{
				// Note: for multi-hop use the server name or hostname instead of the country
				Country:  country,
				Region:   getStringValue(logicalServer.Region),
				City:     getStringValue(logicalServer.City),
				Name:     logicalServer.Name,
				Hostname: physicalServer.Domain,
				EntryIP:  physicalServer.EntryIP,
				ExitIP:   physicalServer.ExitIP,
			}
			servers = append(servers, server)
		}
	}

	sort.Slice(servers, func(i, j int) bool {
		a, b := servers[i], servers[j]
		if a.Country == b.Country { //nolint:nestif
			if a.Region == b.Region {
				if a.City == b.City {
					if a.Name == b.Name {
						return a.Hostname < b.Hostname
					}
					return a.Name < b.Name
				}
				return a.City < b.City
			}
			return a.Region < b.Region
		}
		return a.Country < b.Country
	})

	return servers, warnings, nil
}

func getStringValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

func stringifyProtonvpnServers(servers []models.ProtonvpnServer) (s string) {
	s = "func ProtonvpnServers() []models.ProtonvpnServer {\n"
	s += "	return []models.ProtonvpnServer{\n"
	for _, server := range servers {
		s += "		" + server.String() + ",\n"
	}
	s += "	}\n"
	s += "}"
	return s
}
