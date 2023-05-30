package env

import (
	"github.com/qdm12/gluetun/internal/pprof"
	"github.com/qdm12/gosettings/sources/env"
)

func readPprof() (settings pprof.Settings, err error) {
	settings.Enabled, err = env.BoolPtr("PPROF_ENABLED")
	if err != nil {
		return settings, err
	}

	settings.BlockProfileRate, err = env.IntPtr("PPROF_BLOCK_PROFILE_RATE")
	if err != nil {
		return settings, err
	}

	settings.MutexProfileRate, err = env.IntPtr("PPROF_MUTEX_PROFILE_RATE")
	if err != nil {
		return settings, err
	}

	settings.HTTPServer.Address = env.String("PPROF_HTTP_SERVER_ADDRESS")

	return settings, nil
}
