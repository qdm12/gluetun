package mullvad

import (
	"sort"

	"github.com/qdm12/gluetun/internal/models"
)

func sortServers(servers []models.MullvadServer) {
	sort.Slice(servers, func(i, j int) bool {
		if servers[i].Country == servers[j].Country {
			if servers[i].City == servers[j].City {
				if servers[i].Hostname == servers[j].Hostname {
					return servers[i].ISP < servers[j].ISP
				}
				return servers[i].Hostname < servers[j].Hostname
			}
			return servers[i].City < servers[j].City
		}
		return servers[i].Country < servers[j].Country
	})
}
