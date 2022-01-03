package protonvpn

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

func (p *Protonvpn) GetConnection(selection settings.ServerSelection) (
	connection models.Connection, err error) {
	protocol := constants.UDP
	if *selection.OpenVPN.TCP {
		protocol = constants.TCP
	}

	port, err := getPort(*selection.OpenVPN.TCP, *selection.OpenVPN.CustomPort)
	if err != nil {
		return connection, err
	}

	servers, err := p.filterServers(selection)
	if err != nil {
		return connection, err
	}

	connections := make([]models.Connection, len(servers))
	for i := range servers {
		connections[i] = models.Connection{
			Type:     selection.VPN,
			IP:       servers[i].EntryIP,
			Port:     port,
			Protocol: protocol,
		}
	}

	return utils.PickConnection(connections, selection, p.randSource)
}
