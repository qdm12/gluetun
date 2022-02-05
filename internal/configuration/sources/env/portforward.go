package env

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func (r *Reader) readPortForward() (
	portForwarding settings.PortForwarding, err error) {
	key, _ := r.getEnvWithRetro(
		"PRIVATE_INTERNET_ACCESS_VPN_PORT_FORWARDING",
		"PORT_FORWARDING")
	portForwarding.Enabled, err = envToBoolPtr(key)
	if err != nil {
		return portForwarding, fmt.Errorf("environment variable %s: %w", key, err)
	}

	portForwarding.Filepath = envToStringPtr("PORT_FORWARDING_STATUS_FILE")

	return portForwarding, nil
}
