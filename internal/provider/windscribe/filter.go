package windscribe

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (w *Windscribe) filterServers(selection settings.ServerSelection) (
	servers []models.WindscribeServer, err error) {
	for _, server := range w.servers {
		switch {
		case
			server.VPN != selection.VPN,
			utils.FilterByPossibilities(server.Region, selection.Regions),
			utils.FilterByPossibilities(server.City, selection.Cities),
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
