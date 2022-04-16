package vyprvpn

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (v *Vyprvpn) filterServers(selection settings.ServerSelection) (
	servers []models.Server, err error) {
	for _, server := range v.servers {
		switch {
		case
			utils.FilterByPossibilities(server.Region, selection.Regions),
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
