package ivpn

import (
	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (i *Ivpn) GetConnection(selection configuration.ServerSelection) (
	connection models.Connection, err error) {
	port := getPort(selection)
	protocol := utils.GetProtocol(selection)

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
				PubKey:   server.WgPubKey, // Wireguard only
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
	const (
		defaultOpenVPNTCP = 443
		defaultOpenVPNUDP = 1194
		defaultWireguard  = 58237
	)
	return utils.GetPort(selection, defaultOpenVPNTCP,
		defaultOpenVPNUDP, defaultWireguard)
}
