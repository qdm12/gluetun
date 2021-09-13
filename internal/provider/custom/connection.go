package custom

import (
	"errors"
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

var (
	ErrVPNTypeNotSupported = errors.New("VPN type not supported for custom provider")
	ErrExtractConnection   = errors.New("cannot extract connection")
)

// GetConnection gets the connection from the OpenVPN configuration file.
func (p *Provider) GetConnection(selection configuration.ServerSelection) (
	connection models.Connection, err error) {
	if selection.VPN != constants.OpenVPN {
		return connection, fmt.Errorf("%w: %s", ErrVPNTypeNotSupported, selection.VPN)
	}

	_, connection, err = p.extractor.Data(selection.OpenVPN.ConfFile)
	if err != nil {
		return connection, fmt.Errorf("%w: %s", ErrExtractConnection, err)
	}

	connection.Port = getPort(connection.Port, selection)

	return connection, nil
}

// Port found is overridden by custom port set with `PORT` or `WIREGUARD_PORT`.
func getPort(foundPort uint16, selection configuration.ServerSelection) (port uint16) {
	return utils.GetPort(selection, foundPort, foundPort, foundPort)
}
