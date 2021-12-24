package files

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func (r *Reader) readSystem() (system settings.System, err error) {
	// TODO timezone from /etc/localtime
	return system, nil
}
