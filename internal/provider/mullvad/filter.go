package mullvad

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (m *Mullvad) filterServers(selection settings.ServerSelection) (
	servers []models.Server, err error) {
	for _, server := range m.servers {
		switch {
		case
			server.VPN != selection.VPN,
			utils.FilterByPossibilities(server.Country, selection.Countries),
			utils.FilterByPossibilities(server.City, selection.Cities),
			utils.FilterByPossibilities(server.ISP, selection.ISPs),
			utils.FilterByPossibilities(server.Hostname, selection.Hostnames),
			*selection.OwnedOnly && !server.Owned:
		default:
			servers = append(servers, server)
		}
	}

	if len(servers) == 0 {
		return nil, utils.NoServerFoundError(selection)
	}

	return servers, nil
}
