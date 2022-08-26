package env

import (
	"errors"
	"fmt"
	"net"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/govalid/port"
)

func (s *Source) readWireguardSelection() (
	selection settings.WireguardSelection, err error) {
	selection.EndpointIP, err = s.readWireguardEndpointIP()
	if err != nil {
		return selection, err
	}

	selection.EndpointPort, err = s.readWireguardCustomPort()
	if err != nil {
		return selection, err
	}

	selection.PublicKey = getCleanedEnv("WIREGUARD_PUBLIC_KEY")

	return selection, nil
}

var ErrIPAddressParse = errors.New("cannot parse IP address")

func (s *Source) readWireguardEndpointIP() (endpointIP net.IP, err error) {
	key, value := s.getEnvWithRetro("VPN_ENDPOINT_IP", "WIREGUARD_ENDPOINT_IP")
	if value == "" {
		return nil, nil
	}

	endpointIP = net.ParseIP(value)
	if endpointIP == nil {
		return nil, fmt.Errorf("environment variable %s: %w: %s",
			key, ErrIPAddressParse, value)
	}

	return endpointIP, nil
}

func (s *Source) readWireguardCustomPort() (customPort *uint16, err error) {
	key, value := s.getEnvWithRetro("VPN_ENDPOINT_PORT", "WIREGUARD_ENDPOINT_PORT")
	if value == "" {
		return nil, nil //nolint:nilnil
	}

	customPort = new(uint16)
	*customPort, err = port.Validate(value)
	if err != nil {
		return nil, fmt.Errorf("environment variable %s: %w", key, err)
	}

	return customPort, nil
}
