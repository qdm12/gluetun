package nordvpn

import (
	"sort"

	"github.com/qdm12/gluetun/internal/models"
)

func sortServers(servers []models.NordvpnServer) {
	sort.Slice(servers, func(i, j int) bool {
		if servers[i].Region == servers[j].Region {
			return servers[i].Number < servers[j].Number
		}
		return servers[i].Region < servers[j].Region
	})
}
