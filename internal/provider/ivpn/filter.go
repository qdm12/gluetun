package ivpn

import (
	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (i *Ivpn) filterServers(selection configuration.ServerSelection) (
	servers []models.IvpnServer, err error) {
	for _, server := range i.servers {
		switch {
		case
			server.VPN != selection.VPN,
			utils.FilterByPossibilities(server.ISP, selection.ISPs),
			utils.FilterByPossibilities(server.Country, selection.Countries),
			utils.FilterByPossibilities(server.City, selection.Cities),
			utils.FilterByPossibilities(server.Hostname, selection.Hostnames),
			selection.OpenVPN.TCP && !server.TCP,
			!selection.OpenVPN.TCP && !server.UDP:
		default:
			servers = append(servers, server)
		}
	}

	if len(servers) == 0 {
		return nil, utils.NoServerFoundError(selection)
	}

	return servers, nil
}
