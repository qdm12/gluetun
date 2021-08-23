package ivpn

import (
	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (i *Ivpn) GetConnection(selection configuration.ServerSelection) (
	connection models.Connection, err error) {
	port := getPort(selection)
	protocol := getProtocol(selection.OpenVPN.TCP)

	servers, err := i.filterServers(selection)
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

func getPort(selection configuration.ServerSelection) (port uint16) {
	customPort := selection.OpenVPN.CustomPort
	if customPort > 0 {
		return customPort
	}
	port = 1194
	if selection.OpenVPN.TCP {
		port = 443
	}
	return port
}

func getProtocol(tcp bool) (protocol string) {
	if tcp {
		return constants.TCP
	}
	return constants.UDP
}
