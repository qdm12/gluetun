package protonvpn

import (
	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (p *Protonvpn) GetOpenVPNConnection(selection configuration.ServerSelection) (
	connection models.OpenVPNConnection, err error) {
	protocol := constants.UDP
	if selection.TCP {
		protocol = constants.TCP
	}

	port, err := getPort(selection.TCP, selection.CustomPort)
	if err != nil {
		return connection, err
	}

	servers, err := p.filterServers(selection)
	if err != nil {
		return connection, err
	}

	connections := make([]models.OpenVPNConnection, len(servers))
	for i := range servers {
		connections[i] = models.OpenVPNConnection{
			IP:       servers[i].EntryIP,
			Port:     port,
			Protocol: protocol,
		}
	}

	if selection.TargetIP != nil {
		return utils.GetTargetIPOpenVPNConnection(connections, selection.TargetIP)
	}

	return utils.PickRandomOpenVPNConnection(connections, p.randSource), nil
}
