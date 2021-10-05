package perfectprivacy

import (
	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (p *Perfectprivacy) GetConnection(selection configuration.ServerSelection) (
	connection models.Connection, err error) {
	const defaultPort uint16 = 443
	port := defaultPort
	if selection.OpenVPN.CustomPort > 0 {
		port = selection.OpenVPN.CustomPort
	}
	protocol := utils.GetProtocol(selection)

	servers, err := p.filterServers(selection)
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
