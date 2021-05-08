package surfshark

import (
	"sort"

	"github.com/qdm12/gluetun/internal/models"
)

func sortServers(servers []models.SurfsharkServer) {
	sort.Slice(servers, func(i, j int) bool {
		return servers[i].Region < servers[j].Region
	})
}
