package perfectprivacy

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (p *Perfectprivacy) filterServers(selection settings.ServerSelection) (
	servers []models.PerfectprivacyServer, err error) {
	for _, server := range p.servers {
		switch {
		case
			utils.FilterByPossibilities(server.City, selection.Cities),
			utils.FilterByProtocol(selection, server.TCP, server.UDP):
		default:
			servers = append(servers, server)
		}
	}

	if len(servers) == 0 {
		return nil, utils.NoServerFoundError(selection)
	}

	return servers, nil
}
