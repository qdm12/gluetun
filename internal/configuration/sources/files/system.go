package files

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func (s *Source) readSystem() (system settings.System, err error) {
	// TODO timezone from /etc/localtime
	return system, nil
}
