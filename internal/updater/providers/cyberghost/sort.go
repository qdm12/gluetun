package cyberghost

import (
	"sort"

	"github.com/qdm12/gluetun/internal/models"
)

func sortServers(servers []models.CyberghostServer) {
	sort.Slice(servers, func(i, j int) bool {
		if servers[i].Region == servers[j].Region {
			return servers[i].Group < servers[j].Group
		}
		return servers[i].Region < servers[j].Region
	})
}
