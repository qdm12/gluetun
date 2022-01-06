package cyberghost

import (
	"errors"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

var ErrGroupMismatchesProtocol = errors.New("server group does not match protocol")

func (c *Cyberghost) filterServers(selection settings.ServerSelection) (
	servers []models.CyberghostServer, err error) {
	for _, server := range c.servers {
		switch {
		case
			utils.FilterByProtocol(selection, server.TCP, server.UDP),
			utils.FilterByPossibilities(server.Country, selection.Countries),
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
