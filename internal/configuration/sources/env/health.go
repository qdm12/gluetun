package env

import (
	"fmt"
	"time"

	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func (s *Source) ReadHealth() (health settings.Health, err error) {
	health.ServerAddress = getCleanedEnv("HEALTH_SERVER_ADDRESS")
	_, health.TargetAddress = s.getEnvWithRetro("HEALTH_TARGET_ADDRESS", "HEALTH_ADDRESS_TO_PING")

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
	envKey, value := s.getEnvWithRetro(envKey, retroEnvKey)
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
