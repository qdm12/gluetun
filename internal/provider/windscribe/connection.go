package windscribe

import (
	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (w *Windscribe) GetConnection(selection configuration.ServerSelection) (
	connection models.Connection, err error) {
	protocol := constants.UDP
	var port uint16 = 443
	if selection.OpenVPN.TCP {
		protocol = constants.TCP
		port = 1194
	}

	if selection.OpenVPN.CustomPort > 0 {
		port = selection.OpenVPN.CustomPort
	}

	servers, err := w.filterServers(selection)
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
				Hostname: server.OvpnX509,
			}
			connections = append(connections, connection)
		}
	}

	if selection.TargetIP != nil {
		return utils.GetTargetIPConnection(connections, selection.TargetIP)
	}

	return utils.PickRandomConnection(connections, w.randSource), nil
}
