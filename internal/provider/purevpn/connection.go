package purevpn

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (p *Purevpn) GetConnection(selection settings.ServerSelection) (
	connection models.Connection, err error) {
	protocol := constants.UDP
	var port uint16 = 53
	if *selection.OpenVPN.TCP {
		protocol = constants.TCP
		port = 80
	}

	servers, err := utils.FilterServers(p.servers, selection)
	if err != nil {
		return connection, err
	}

	var connections []models.Connection
	for _, server := range servers {
		for _, IP := range server.IPs {
			connection := models.Connection{
				Type:     selection.VPN,
				IP:       IP,
				Port:     port,
				Protocol: protocol,
			}
			connections = append(connections, connection)
		}
	}

	return utils.PickConnection(connections, selection, p.randSource)
}
