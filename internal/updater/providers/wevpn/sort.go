package wevpn

import (
	"sort"

	"github.com/qdm12/gluetun/internal/models"
)

func sortServers(servers []models.WevpnServer) {
	sort.Slice(servers, func(i, j int) bool {
		if servers[i].City == servers[j].City {
			return servers[i].Hostname < servers[j].Hostname
		}
		return servers[i].City < servers[j].City
	})
}
