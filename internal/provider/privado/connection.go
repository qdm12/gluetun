package privado

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (p *Privado) GetConnection(selection settings.ServerSelection) (
	connection models.Connection, err error) {
	const port = 1194
	const protocol = constants.UDP

	servers, err := utils.FilterServers(p.servers, selection)
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
				Hostname: server.Hostname,
			}
			connections = append(connections, connection)
		}
	}

	return utils.PickConnection(connections, selection, p.randSource)
}
