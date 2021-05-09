package pia

import (
	"sort"

	"github.com/qdm12/gluetun/internal/models"
)

func sortServers(servers []models.PIAServer) {
	sort.Slice(servers, func(i, j int) bool {
		if servers[i].Region == servers[j].Region {
			if servers[i].Hostname == servers[j].Hostname {
				return servers[i].ServerName < servers[j].ServerName
			}
			return servers[i].Hostname < servers[j].Hostname
		}
		return servers[i].Region < servers[j].Region
	})
}
