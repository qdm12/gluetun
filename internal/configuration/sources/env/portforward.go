package env

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func readPortForward() (
	portForwarding settings.PortForwarding, err error) {
	portForwarding.Enabled, err = envToBoolPtr("PORT_FORWARDING")
	if err != nil {
		return portForwarding, fmt.Errorf("environment variable PORT_FORWARDING: %w", err)
	}

	portForwarding.Filepath = envToStringPtr("PORT_FORWARDING_STATUS_FILE")

	return portForwarding, nil
}
