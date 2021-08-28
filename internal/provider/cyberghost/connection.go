package cyberghost

import (
	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (c *Cyberghost) GetConnection(selection configuration.ServerSelection) (
	connection models.Connection, err error) {
	const port = 443
	protocol := constants.UDP
	if selection.OpenVPN.TCP {
		protocol = constants.TCP
	}

	servers, err := c.filterServers(selection)
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

	return utils.PickConnection(connections, selection, c.randSource)
}
