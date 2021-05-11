package cyberghost

import (
	"strings"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (c *Cyberghost) filterServers(selection configuration.ServerSelection) (
	servers []models.CyberghostServer, err error) {
	for _, server := range c.servers {
		switch {
		case selection.Group != "" && !strings.EqualFold(selection.Group, server.Group), // TODO make CSV
			utils.FilterByPossibilities(server.Region, selection.Regions),
			utils.FilterByPossibilities(server.Hostname, selection.Hostnames):
		default:
			servers = append(servers, server)
		}
	}

	if len(servers) == 0 {
		return nil, utils.NoServerFoundError(selection)
	}

	return servers, nil
}
