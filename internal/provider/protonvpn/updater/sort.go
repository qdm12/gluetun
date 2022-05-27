package protonvpn

import (
	"sort"

	"github.com/qdm12/gluetun/internal/models"
)

func sortServers(servers []models.Server) {
	sort.Slice(servers, func(i, j int) bool {
		a, b := servers[i], servers[j]
		if a.Country == b.Country { //nolint:nestif
			if a.Region == b.Region {
				if a.City == b.City {
					if a.ServerName == b.ServerName {
						return a.Hostname < b.Hostname
					}
					return a.ServerName < b.ServerName
				}
				return a.City < b.City
			}
			return a.Region < b.Region
		}
		return a.Country < b.Country
	})
}
