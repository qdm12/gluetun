package vyprvpn

import (
	"sort"

	"github.com/qdm12/gluetun/internal/models"
)

func sortServers(servers []models.VyprvpnServer) {
	sort.Slice(servers, func(i, j int) bool {
		return servers[i].Region < servers[j].Region
	})
}
