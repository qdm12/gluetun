package privado

import (
	"errors"
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

var ErrProtocolUnsupported = errors.New("network protocol is not supported")

func (p *Privado) GetConnection(selection settings.ServerSelection) (
	connection models.Connection, err error) {
	const port = 1194
	const protocol = constants.UDP
	if *selection.OpenVPN.TCP {
		return connection, fmt.Errorf("%w: TCP for provider Privado", ErrProtocolUnsupported)
	}

	servers, err := p.filterServers(selection)
	if err != nil {
		return connection, err
	}

	connections := make([]models.Connection, len(servers))
	for i := range servers {
		connection := models.Connection{
			Type:     selection.VPN,
			IP:       servers[i].IP,
			Port:     port,
			Protocol: protocol,
			Hostname: servers[i].Hostname,
		}
		connections[i] = connection
	}

	return utils.PickConnection(connections, selection, p.randSource)
}
