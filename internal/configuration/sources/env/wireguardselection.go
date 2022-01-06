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
	selection.EndpointIP, err = readWireguardEndpointIP()
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

func readWireguardEndpointIP() (endpointIP net.IP, err error) {
	s := os.Getenv("WIREGUARD_ENDPOINT_IP")
	if s == "" {
		return nil, nil
	}
	endpointIP = net.ParseIP(s)
	if endpointIP == nil {
		return nil, fmt.Errorf("environment variable WIREGUARD_ENDPOINT_IP: %w: %s",
			ErrIPAddressParse, s)
	}
	return endpointIP, nil
}

func (r *Reader) readWireguardCustomPort() (customPort *uint16, err error) {
	key := "WIREGUARD_ENDPOINT_PORT"
	s := os.Getenv(key)
	if s == "" {
		// Retro-compatibility
		key = "WIREGUARD_PORT"
		s = os.Getenv(key)
		if s == "" {
			return nil, nil //nolint:nilnil
		}
		r.onRetroActive("WIREGUARD_PORT", "WIREGUARD_ENDPOINT_PORT")
	}

	customPort = new(uint16)
	*customPort, err = port.Validate(s)
	if err != nil {
		return nil, fmt.Errorf("environment variable %s: %w", key, err)
	}

	return customPort, nil
}
