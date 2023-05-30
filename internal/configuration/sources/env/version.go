package env

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gosettings/sources/env"
)

func readVersion() (version settings.Version, err error) {
	version.Enabled, err = env.BoolPtr("VERSION_INFORMATION")
	if err != nil {
		return version, err
	}

	return version, nil
}
