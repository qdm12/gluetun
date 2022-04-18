package privatevpn

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (p *Privatevpn) GetConnection(selection settings.ServerSelection) (
	connection models.Connection, err error) {
	protocol := constants.UDP
	var port uint16 = 1194
	if *selection.OpenVPN.TCP {
		protocol = constants.TCP
		port = 443
	}
	if *selection.OpenVPN.CustomPort > 0 {
		port = *selection.OpenVPN.CustomPort
	}

	servers, err := utils.FilterServers(p.servers, selection)
	if err != nil {
		return connection, err
	}

	var connections []models.Connection
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

	return utils.PickConnection(connections, selection, p.randSource)
}
