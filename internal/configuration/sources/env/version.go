package env

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func (s *Source) readVersion() (version settings.Version, err error) {
	version.Enabled, err = s.env.BoolPtr("VERSION_INFORMATION")
	if err != nil {
		return version, err
	}

	return version, nil
}
