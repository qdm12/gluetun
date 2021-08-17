package vyprvpn

import (
	"errors"
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

var ErrProtocolUnsupported = errors.New("network protocol is not supported")

func (v *Vyprvpn) GetOpenVPNConnection(selection configuration.ServerSelection) (
	connection models.OpenVPNConnection, err error) {
	const port = 443
	const protocol = constants.UDP
	if selection.OpenVPN.TCP {
		return connection, fmt.Errorf("%w: TCP for provider VyprVPN", ErrProtocolUnsupported)
	}

	servers, err := v.filterServers(selection)
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
		return utils.GetTargetIPOpenVPNConnection(connections, selection.TargetIP)
	}

	return utils.PickRandomOpenVPNConnection(connections, v.randSource), nil
}
