package env

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gosettings/sources/env"
)

func (s *Source) readSystem() (system settings.System, err error) {
	system.PUID, err = s.env.Uint32Ptr("PUID", env.RetroKeys("UID"))
	if err != nil {
		return system, err
	}

	system.PGID, err = s.env.Uint32Ptr("PGID", env.RetroKeys("GID"))
	if err != nil {
		return system, err
	}

	system.Timezone = s.env.String("TZ")

	return system, nil
}
