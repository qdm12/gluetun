package expressvpn

import (
	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (p *Provider) filterServers(selection configuration.ServerSelection) (
	servers []models.ExpressvpnServer, err error) {
	for _, server := range p.servers {
		switch {
		case
			utils.FilterByPossibilities(server.Country, selection.Countries),
			utils.FilterByPossibilities(server.City, selection.Cities),
			utils.FilterByPossibilities(server.Hostname, selection.Hostnames),
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
