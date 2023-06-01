package env

import (
	"github.com/qdm12/gluetun/internal/pprof"
)

func (s *Source) readPprof() (settings pprof.Settings, err error) {
	settings.Enabled, err = s.env.BoolPtr("PPROF_ENABLED")
	if err != nil {
		return settings, err
	}

	settings.BlockProfileRate, err = s.env.IntPtr("PPROF_BLOCK_PROFILE_RATE")
	if err != nil {
		return settings, err
	}

	settings.MutexProfileRate, err = s.env.IntPtr("PPROF_MUTEX_PROFILE_RATE")
	if err != nil {
		return settings, err
	}

	settings.HTTPServer.Address = s.env.String("PPROF_HTTP_SERVER_ADDRESS")

	return settings, nil
}
