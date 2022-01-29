package env

import (
	"fmt"
	"os"
	"time"

	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func (r *Reader) ReadHealth() (health settings.Health, err error) {
	health.ServerAddress = os.Getenv("HEALTH_SERVER_ADDRESS")
	health.AddressToPing = os.Getenv("HEALTH_ADDRESS_TO_PING")

	health.VPN.Initial, err = r.readDurationWithRetro(
		"HEALTH_VPN_DURATION_INITIAL",
		"HEALTH_OPENVPN_DURATION_INITIAL")
	if err != nil {
		return health, err
	}

	health.VPN.Initial, err = r.readDurationWithRetro(
		"HEALTH_VPN_DURATION_ADDITION",
		"HEALTH_OPENVPN_DURATION_ADDITION")
	if err != nil {
		return health, err
	}

	return health, nil
}

func (r *Reader) readDurationWithRetro(envKey, retroEnvKey string) (d *time.Duration, err error) {
	envKey, s := r.getEnvWithRetro(envKey, retroEnvKey)
	if s == "" {
		return nil, nil //nolint:nilnil
	}

	d = new(time.Duration)
	*d, err = time.ParseDuration(s)
	if err != nil {
		return nil, fmt.Errorf(
			"environment variable %s: %w",
			envKey, err)
	}

	return d, nil
}
