package perfectprivacy

import (
	"sort"

	"github.com/qdm12/gluetun/internal/models"
)

func sortServers(servers []models.Server) {
	sort.Slice(servers, func(i, j int) bool {
		return servers[i].City < servers[j].City
	})
}
