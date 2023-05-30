package env

import (
	"fmt"
	"net/netip"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gosettings/sources/env"
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

	selection.PublicKey = env.Get("WIREGUARD_PUBLIC_KEY", env.ForceLowercase(false))

	return selection, nil
}

func (s *Source) readWireguardEndpointIP() (endpointIP netip.Addr, err error) {
	key, value := s.getEnvWithRetro("VPN_ENDPOINT_IP", []string{"WIREGUARD_ENDPOINT_IP"})
	if value == "" {
		return endpointIP, nil
	}

	endpointIP, err = netip.ParseAddr(value)
	if err != nil {
		return endpointIP, fmt.Errorf("environment variable %s: %w", key, err)
	}

	return endpointIP, nil
}

func (s *Source) readWireguardCustomPort() (customPort *uint16, err error) {
	key, value := s.getEnvWithRetro("VPN_ENDPOINT_PORT", []string{"WIREGUARD_ENDPOINT_PORT"})
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
