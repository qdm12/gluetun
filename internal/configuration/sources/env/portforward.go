package env

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gosettings/sources/env"
)

func (s *Source) readPortForward() (
	portForwarding settings.PortForwarding, err error) {
	key, _ := s.getEnvWithRetro("VPN_PORT_FORWARDING",
		[]string{
			"PRIVATE_INTERNET_ACCESS_VPN_PORT_FORWARDING",
			"PORT_FORWARDING",
		})
	portForwarding.Enabled, err = s.env.BoolPtr(key)
	if err != nil {
		return portForwarding, err
	}

	portForwarding.Provider = s.env.Get("VPN_PORT_FORWARDING_PROVIDER")

	_, portForwarding.Filepath = s.getEnvWithRetro("VPN_PORT_FORWARDING_STATUS_FILE",
		[]string{
			"VPN_PORT_FORWARDING_STATUS_FILE",
			"PRIVATE_INTERNET_ACCESS_VPN_PORT_FORWARDING_STATUS_FILE",
			"PORT_FORWARDING_STATUS_FILE",
		}, env.ForceLowercase(false))

	return portForwarding, nil
}
