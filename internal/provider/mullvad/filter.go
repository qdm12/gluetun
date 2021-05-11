package mullvad

import (
	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (m *Mullvad) filterServers(selection configuration.ServerSelection) (
	servers []models.MullvadServer, err error) {
	for _, server := range m.servers {
		switch {
		case
			utils.FilterByPossibilities(server.Country, selection.Countries),
			utils.FilterByPossibilities(server.City, selection.Cities),
			utils.FilterByPossibilities(server.ISP, selection.ISPs),
			utils.FilterByPossibilities(server.Hostname, selection.Hostnames),
			selection.Owned && !server.Owned:
		default:
			servers = append(servers, server)
		}
	}

	if len(servers) == 0 {
		return nil, utils.NoServerFoundError(selection)
	}

	return servers, nil
}
