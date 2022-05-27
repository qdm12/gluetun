package windscribe

import (
	"sort"

	"github.com/qdm12/gluetun/internal/models"
)

func sortServers(servers []models.Server) {
	sort.Slice(servers, func(i, j int) bool {
		if servers[i].Region == servers[j].Region {
			if servers[i].City == servers[j].City {
				if servers[i].Hostname == servers[j].Hostname {
					return servers[i].VPN < servers[j].VPN
				}
				return servers[i].Hostname < servers[j].Hostname
			}
			return servers[i].City < servers[j].City
		}
		return servers[i].Region < servers[j].Region
	})
}
