package nordvpn

import (
	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (n *Nordvpn) GetOpenVPNConnection(selection configuration.ServerSelection) (
	connection models.OpenVPNConnection, err error) {
	var port uint16 = 1194
	protocol := constants.UDP
	if selection.OpenVPN.TCP {
		port = 443
		protocol = constants.TCP
	}

	servers, err := n.filterServers(selection)
	if err != nil {
		return connection, err
	}

	connections := make([]models.OpenVPNConnection, len(servers))
	for i := range servers {
		connection := models.OpenVPNConnection{
			IP:       servers[i].IP,
			Port:     port,
			Protocol: protocol,
		}
		connections[i] = connection
	}

	if selection.TargetIP != nil {
		return utils.GetTargetIPOpenVPNConnection(connections, selection.TargetIP)
	}

	return utils.PickRandomOpenVPNConnection(connections, n.randSource), nil
}
