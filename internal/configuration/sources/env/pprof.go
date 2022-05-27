package env

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/pprof"
)

func readPprof() (settings pprof.Settings, err error) {
	settings.Enabled, err = envToBoolPtr("PPROF_ENABLED")
	if err != nil {
		return settings, fmt.Errorf("environment variable PPROF_ENABLED: %w", err)
	}

	settings.BlockProfileRate, err = envToInt("PPROF_BLOCK_PROFILE_RATE")
	if err != nil {
		return settings, fmt.Errorf("environment variable PPROF_BLOCK_PROFILE_RATE: %w", err)
	}

	settings.MutexProfileRate, err = envToInt("PPROF_MUTEX_PROFILE_RATE")
	if err != nil {
		return settings, fmt.Errorf("environment variable PPROF_MUTEX_PROFILE_RATE: %w", err)
	}

	settings.HTTPServer.Address = getCleanedEnv("PPROF_HTTP_SERVER_ADDRESS")

	return settings, nil
}
