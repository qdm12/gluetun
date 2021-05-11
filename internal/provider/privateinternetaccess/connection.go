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
	if selection.TCP {
		protocol = constants.TCP
	}

	port, err := getPort(selection.TCP, selection.EncryptionPreset, selection.CustomPort)
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
			}
			connections = append(connections, connection)
		}
	}

	if selection.TargetIP != nil {
		connection, err = utils.GetTargetIPConnection(connections, selection.TargetIP)
	} else {
		connection, err = utils.PickRandomConnection(connections, p.randSource), nil
	}

	if err != nil {
		return connection, err
	}

	p.activeServer = findActiveServer(servers, connection)

	return connection, nil
}

func findActiveServer(servers []models.PIAServer,
	connection models.OpenVPNConnection) (activeServer models.PIAServer) {
	// Reverse lookup server using the randomly picked connection
	for _, server := range servers {
		for _, ip := range server.IPs {
			if connection.IP.Equal(ip) {
				return server
			}
		}
	}
	return activeServer
}
