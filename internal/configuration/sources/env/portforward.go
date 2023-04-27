package env

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func (s *Source) readPortForward() (
	portForwarding settings.PortForwarding, err error) {
	key, _ := s.getEnvWithRetro(
		"VPN_PORT_FORWARDING",
		"PRIVATE_INTERNET_ACCESS_VPN_PORT_FORWARDING",
		"PORT_FORWARDING")
	portForwarding.Enabled, err = envToBoolPtr(key)
	if err != nil {
		return portForwarding, fmt.Errorf("environment variable %s: %w", key, err)
	}

	_, value := s.getEnvWithRetro(
		"VPN_PORT_FORWARDING_STATUS_FILE",
		"PRIVATE_INTERNET_ACCESS_VPN_PORT_FORWARDING_STATUS_FILE",
		"PORT_FORWARDING_STATUS_FILE")
	if value != "" {
		portForwarding.Filepath = stringPtr(value)
	}

	return portForwarding, nil
}
