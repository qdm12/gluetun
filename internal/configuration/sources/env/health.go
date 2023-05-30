package env

import (
	"time"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gosettings/sources/env"
)

func (s *Source) ReadHealth() (health settings.Health, err error) {
	health.ServerAddress = env.String("HEALTH_SERVER_ADDRESS")
	targetAddressEnvKey, _ := s.getEnvWithRetro("HEALTH_TARGET_ADDRESS", []string{"HEALTH_ADDRESS_TO_PING"})
	health.TargetAddress = env.String(targetAddressEnvKey)

	successWaitPtr, err := env.DurationPtr("HEALTH_SUCCESS_WAIT_DURATION")
	if err != nil {
		return health, err
	} else if successWaitPtr != nil {
		health.SuccessWait = *successWaitPtr
	}

	health.VPN.Initial, err = s.readDurationWithRetro(
		"HEALTH_VPN_DURATION_INITIAL",
		"HEALTH_OPENVPN_DURATION_INITIAL")
	if err != nil {
		return health, err
	}

	health.VPN.Addition, err = s.readDurationWithRetro(
		"HEALTH_VPN_DURATION_ADDITION",
		"HEALTH_OPENVPN_DURATION_ADDITION")
	if err != nil {
		return health, err
	}

	return health, nil
}

func (s *Source) readDurationWithRetro(envKey, retroEnvKey string) (d *time.Duration, err error) {
	envKey, _ = s.getEnvWithRetro(envKey, []string{retroEnvKey})
	return env.DurationPtr(envKey)
}
