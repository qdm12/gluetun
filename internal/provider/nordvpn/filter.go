package nordvpn

import (
	"strconv"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (n *Nordvpn) filterServers(selection settings.ServerSelection) (
	servers []models.NordvpnServer, err error) {
	selectedNumbers := make([]string, len(selection.Numbers))
	for i := range selection.Numbers {
		selectedNumbers[i] = strconv.Itoa(int(selection.Numbers[i]))
	}

	for _, server := range n.servers {
		serverNumber := strconv.Itoa(int(server.Number))
		switch {
		case
			utils.FilterByPossibilities(server.Region, selection.Regions),
			utils.FilterByPossibilities(server.Hostname, selection.Hostnames),
			utils.FilterByPossibilities(serverNumber, selectedNumbers),
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
