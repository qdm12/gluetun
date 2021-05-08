package purevpn

import (
	"sort"

	"github.com/qdm12/gluetun/internal/models"
)

func sortServers(servers []models.PurevpnServer) {
	sort.Slice(servers, func(i, j int) bool {
		if servers[i].Country == servers[j].Country {
			if servers[i].Region == servers[j].Region {
				return servers[i].City < servers[j].City
			}
			return servers[i].Region < servers[j].Region
		}
		return servers[i].Country < servers[j].Country
	})
}
