package env

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/pprof"
	"github.com/qdm12/gosettings/sources/env"
)

func readPprof() (settings pprof.Settings, err error) {
	settings.Enabled, err = env.BoolPtr("PPROF_ENABLED")
	if err != nil {
		return settings, fmt.Errorf("environment variable PPROF_ENABLED: %w", err)
	}

	settings.BlockProfileRate, err = env.IntPtr("PPROF_BLOCK_PROFILE_RATE")
	if err != nil {
		return settings, fmt.Errorf("environment variable PPROF_BLOCK_PROFILE_RATE: %w", err)
	}

	settings.MutexProfileRate, err = env.IntPtr("PPROF_MUTEX_PROFILE_RATE")
	if err != nil {
		return settings, fmt.Errorf("environment variable PPROF_MUTEX_PROFILE_RATE: %w", err)
	}

	settings.HTTPServer.Address = env.Get("PPROF_HTTP_SERVER_ADDRESS")

	return settings, nil
}
