package env

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gosettings/sources/env"
)

func (s *Source) readPortForward() (
	portForwarding settings.PortForwarding, err error) {
	portForwarding.Enabled, err = s.env.BoolPtr("VPN_PORT_FORWARDING",
		env.RetroKeys(
			"PORT_FORWARDING",
			"PRIVATE_INTERNET_ACCESS_VPN_PORT_FORWARDING",
		))
	if err != nil {
		return portForwarding, err
	}

	portForwarding.Provider = s.env.Get("VPN_PORT_FORWARDING_PROVIDER")

	portForwarding.Filepath = s.env.Get("VPN_PORT_FORWARDING_STATUS_FILE",
		env.ForceLowercase(false),
		env.RetroKeys(
			"PORT_FORWARDING_STATUS_FILE",
			"PRIVATE_INTERNET_ACCESS_VPN_PORT_FORWARDING_STATUS_FILE",
		))

	return portForwarding, nil
}
