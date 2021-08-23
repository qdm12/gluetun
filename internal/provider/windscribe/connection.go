package windscribe

import (
	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (w *Windscribe) GetConnection(selection configuration.ServerSelection) (
	connection models.Connection, err error) {
	port := getPort(selection)
	protocol := utils.GetProtocol(selection)

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
				PubKey:   server.WgPubKey,
			}
			connections = append(connections, connection)
		}
	}

	if selection.TargetIP != nil {
		return utils.GetTargetIPConnection(connections, selection.TargetIP)
	}

	return utils.PickRandomConnection(connections, w.randSource), nil
}

func getPort(selection configuration.ServerSelection) (port uint16) {
	const (
		defaultOpenVPNTCP = 443
		defaultOpenVPNUDP = 1194
		defaultWireguard  = 1194
	)
	return utils.GetPort(selection, defaultOpenVPNTCP,
		defaultOpenVPNUDP, defaultWireguard)
}
