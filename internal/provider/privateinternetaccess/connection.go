package privateinternetaccess

import (
	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (p *PIA) GetOpenVPNConnection(selection configuration.ServerSelection) (
	connection models.OpenVPNConnection, err error) {
	protocol := constants.UDP
	if selection.OpenVPN.TCP {
		protocol = constants.TCP
	}

	port, err := getPort(selection.OpenVPN)
	if err != nil {
		return connection, err
	}

	servers, err := p.filterServers(selection)
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
				Hostname: server.ServerName, // used for port forwarding TLS
			}
			connections = append(connections, connection)
		}
	}

	if selection.TargetIP != nil {
		connection, err = utils.GetTargetIPOpenVPNConnection(connections, selection.TargetIP)
	} else {
		connection, err = utils.PickRandomOpenVPNConnection(connections, p.randSource), nil
	}

	if err != nil {
		return connection, err
	}

	return connection, nil
}
