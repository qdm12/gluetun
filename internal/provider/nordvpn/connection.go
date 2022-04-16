package nordvpn

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (n *Nordvpn) GetConnection(selection settings.ServerSelection) (
	connection models.Connection, err error) {
	var port uint16 = 1194
	protocol := constants.UDP
	if *selection.OpenVPN.TCP {
		port = 443
		protocol = constants.TCP
	}

	servers, err := n.filterServers(selection)
	if err != nil {
		return connection, err
	}

	connections := make([]models.Connection, 0, len(servers))
	for _, server := range servers {
		for _, ip := range server.IPs {
			connection := models.Connection{
				Type:     selection.VPN,
				IP:       ip,
				Port:     port,
				Protocol: protocol,
			}
			connections = append(connections, connection)
		}
	}

	return utils.PickConnection(connections, selection, n.randSource)
}
