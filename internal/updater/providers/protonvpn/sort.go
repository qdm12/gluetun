package protonvpn

import (
	"sort"

	"github.com/qdm12/gluetun/internal/models"
)

func sortServers(servers []models.ProtonvpnServer) {
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
}
