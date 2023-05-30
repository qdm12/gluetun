package env

import (
	"fmt"
	"time"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gosettings/sources/env"
)

func (s *Source) ReadHealth() (health settings.Health, err error) {
	health.ServerAddress = env.Get("HEALTH_SERVER_ADDRESS")
	_, health.TargetAddress = s.getEnvWithRetro("HEALTH_TARGET_ADDRESS", []string{"HEALTH_ADDRESS_TO_PING"})

	successWaitPtr, err := envToDurationPtr("HEALTH_SUCCESS_WAIT_DURATION")
	if err != nil {
		return health, fmt.Errorf("environment variable HEALTH_SUCCESS_WAIT_DURATION: %w", err)
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
	envKey, value := s.getEnvWithRetro(envKey, []string{retroEnvKey})
	if value == "" {
		return nil, nil //nolint:nilnil
	}

	d = new(time.Duration)
	*d, err = time.ParseDuration(value)
	if err != nil {
		return nil, fmt.Errorf("environment variable %s: %w", envKey, err)
	}

	return d, nil
}
