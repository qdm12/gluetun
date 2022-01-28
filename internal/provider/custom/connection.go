package custom

import (
	"errors"
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/openvpn/extract"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

var (
	ErrVPNTypeNotSupported = errors.New("VPN type not supported for custom provider")
	ErrExtractConnection   = errors.New("cannot extract connection")
)

// GetConnection gets the connection from the OpenVPN configuration file.
func (p *Provider) GetConnection(selection settings.ServerSelection) (
	connection models.Connection, err error) {
	switch selection.VPN {
	case constants.OpenVPN:
		return getOpenVPNConnection(p.extractor, selection)
	case constants.Wireguard:
		return getWireguardConnection(selection), nil
	default:
		return connection, fmt.Errorf("%w: %s", ErrVPNTypeNotSupported, selection.VPN)
	}
}

func getOpenVPNConnection(extractor extract.Interface,
	selection settings.ServerSelection) (
	connection models.Connection, err error) {
	_, connection, err = extractor.Data(*selection.OpenVPN.ConfFile)
	if err != nil {
		return connection, fmt.Errorf("%w: %s", ErrExtractConnection, err)
	}

	connection.Port = getPort(connection.Port, selection)
	return connection, nil
}

func getWireguardConnection(selection settings.ServerSelection) (
	connection models.Connection) {
	port := getPort(*selection.Wireguard.EndpointPort, selection)
	return models.Connection{
		Type:     constants.Wireguard,
		IP:       selection.Wireguard.EndpointIP,
		Port:     port,
		Protocol: constants.UDP,
		PubKey:   selection.Wireguard.PublicKey,
	}
}

// Port found is overridden by custom port set with `VPN_ENDPOINT_PORT`.
func getPort(foundPort uint16, selection settings.ServerSelection) (port uint16) {
	return utils.GetPort(selection, foundPort, foundPort, foundPort)
}
