package mullvad

import (
	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (m *Mullvad) GetConnection(selection configuration.ServerSelection) (
	connection models.Connection, err error) {
	var port uint16 = 1194
	protocol := constants.UDP
	if selection.OpenVPN.TCP {
		port = 443
		protocol = constants.TCP
	}

	if selection.OpenVPN.CustomPort > 0 {
		port = selection.OpenVPN.CustomPort
	}

	servers, err := m.filterServers(selection)
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

	if selection.TargetIP != nil {
		return utils.GetTargetIPConnection(connections, selection.TargetIP)
	}

	return utils.PickRandomConnection(connections, m.randSource), nil
}
