package surfshark

import (
	"sort"

	"github.com/qdm12/gluetun/internal/models"
)

func sortServers(servers []models.SurfsharkServer) {
	sort.Slice(servers, func(i, j int) bool {
		if servers[i].Region == servers[j].Region {
			return servers[i].Hostname < servers[j].Hostname
		}
		return servers[i].Region < servers[j].Region
	})
}
