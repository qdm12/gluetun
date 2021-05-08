package fastestvpn

import (
	"sort"

	"github.com/qdm12/gluetun/internal/models"
)

func sortServers(servers []models.FastestvpnServer) {
	sort.Slice(servers, func(i, j int) bool {
		if servers[i].Country == servers[j].Country {
			return servers[i].Hostname < servers[j].Hostname
		}
		return servers[i].Country < servers[j].Country
	})
}
