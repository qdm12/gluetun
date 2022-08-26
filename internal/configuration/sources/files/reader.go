package files

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
)

type Reader struct{}

func New() *Reader {
	return &Reader{}
}

func (r *Reader) String() string { return "files" }

func (r *Reader) Read() (settings settings.Settings, err error) {
	settings.VPN, err = r.readVPN()
	if err != nil {
		return settings, err
	}

	settings.System, err = r.readSystem()
	if err != nil {
		return settings, err
	}

	return settings, nil
}
