package env

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gosettings/sources/env"
)

func (s *Source) ReadHealth() (health settings.Health, err error) {
	health.ServerAddress = s.env.String("HEALTH_SERVER_ADDRESS")
	health.TargetAddress = s.env.String("HEALTH_TARGET_ADDRESS",
		env.RetroKeys("HEALTH_ADDRESS_TO_PING"))

	successWaitPtr, err := s.env.DurationPtr("HEALTH_SUCCESS_WAIT_DURATION")
	if err != nil {
		return health, err
	} else if successWaitPtr != nil {
		health.SuccessWait = *successWaitPtr
	}

	health.VPN.Initial, err = s.env.DurationPtr(
		"HEALTH_VPN_DURATION_INITIAL",
		env.RetroKeys("HEALTH_OPENVPN_DURATION_INITIAL"))
	if err != nil {
		return health, err
	}

	health.VPN.Addition, err = s.env.DurationPtr(
		"HEALTH_VPN_DURATION_ADDITION",
		env.RetroKeys("HEALTH_OPENVPN_DURATION_ADDITION"))
	if err != nil {
		return health, err
	}

	return health, nil
}
