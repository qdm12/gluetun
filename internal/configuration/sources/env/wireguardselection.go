package env

import (
	"errors"
	"fmt"
	"net"
	"os"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/govalid/port"
)

func (r *Reader) readWireguardSelection() (
	selection settings.WireguardSelection, err error) {
	selection.EndpointIP, err = r.readWireguardEndpointIP()
	if err != nil {
		return selection, err
	}

	selection.EndpointPort, err = r.readWireguardCustomPort()
	if err != nil {
		return selection, err
	}

	selection.PublicKey = os.Getenv("WIREGUARD_PUBLIC_KEY")

	return selection, nil
}

var ErrIPAddressParse = errors.New("cannot parse IP address")

func (r *Reader) readWireguardEndpointIP() (endpointIP net.IP, err error) {
	const currentKey = "VPN_ENDPOINT_IP"
	key := "WIREGUARD_ENDPOINT_IP"
	s := os.Getenv(key) // Retro-compatibility
	if s == "" {
		key = currentKey
		s = os.Getenv(key)
		if s == "" {
			return nil, nil
		}
	}

	if key != currentKey {
		r.onRetroActive(key, currentKey)
	}

	endpointIP = net.ParseIP(s)
	if endpointIP == nil {
		return nil, fmt.Errorf("environment variable %s: %w: %s",
			key, ErrIPAddressParse, s)
	}

	return endpointIP, nil
}

func (r *Reader) readWireguardCustomPort() (customPort *uint16, err error) {
	const currentKey = "VPN_ENDPOINT_PORT"
	key := "WIREGUARD_PORT" // Retro-compatibility
	s := os.Getenv(key)
	if s == "" {
		key = "WIREGUARD_ENDPOINT_PORT" // Retro-compatibility
		s = os.Getenv(key)
		if s == "" {
			key = currentKey
			s = os.Getenv(key)
			if s == "" {
				return nil, nil //nolint:nilnil
			}
		}
	}

	if key != currentKey {
		r.onRetroActive(key, currentKey)
	}

	customPort = new(uint16)
	*customPort, err = port.Validate(s)
	if err != nil {
		return nil, fmt.Errorf("environment variable %s: %w", key, err)
	}

	return customPort, nil
}
