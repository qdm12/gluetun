package fastestvpn

import (
	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (f *Fastestvpn) filterServers(selection configuration.ServerSelection) (
	servers []models.FastestvpnServer, err error) {
	for _, server := range f.servers {
		switch {
		case
			utils.FilterByPossibilities(server.Country, selection.Countries),
			utils.FilterByPossibilities(server.Hostname, selection.Hostnames),
			selection.TCP && !server.TCP,
			!selection.TCP && !server.UDP:
		default:
			servers = append(servers, server)
		}
	}

	if len(servers) == 0 {
		return nil, utils.NoServerFoundError(selection)
	}

	return servers, nil
}
