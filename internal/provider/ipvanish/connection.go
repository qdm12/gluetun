package ipvanish

import (
	"errors"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

var ErrProtocolUnsupported = errors.New("network protocol is not supported")

func (i *Ipvanish) GetOpenVPNConnection(selection configuration.ServerSelection) (
	connection models.OpenVPNConnection, err error) {
	const port = 443
	const protocol = constants.UDP
	if selection.TCP {
		return connection, ErrProtocolUnsupported
	}

	servers, err := i.filterServers(selection)
	if err != nil {
		return connection, err
	}

	var connections []models.OpenVPNConnection
	for _, server := range servers {
		for _, IP := range server.IPs {
			connection := models.OpenVPNConnection{
				IP:       IP,
				Port:     port,
				Protocol: protocol,
				Hostname: server.Hostname,
			}
			connections = append(connections, connection)
		}
	}

	if selection.TargetIP != nil {
		return utils.GetTargetIPConnection(connections, selection.TargetIP)
	}

	return utils.PickRandomConnection(connections, i.randSource), nil
}
