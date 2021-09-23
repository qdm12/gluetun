package cyberghost

import (
	"sort"

	"github.com/qdm12/gluetun/internal/models"
)

func sortServers(servers []models.CyberghostServer) {
	sort.Slice(servers, func(i, j int) bool {
		if servers[i].Country == servers[j].Country {
			if servers[i].Group == servers[j].Group {
				return servers[i].Hostname < servers[j].Hostname
			}
			return servers[i].Group < servers[j].Group
		}
		return servers[i].Country < servers[j].Country
	})
}
