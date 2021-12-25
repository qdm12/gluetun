package secrets

import "github.com/qdm12/gluetun/internal/configuration/settings"

type Reader struct {
}

func New() *Reader {
	return &Reader{}
}

func (r *Reader) Read() (settings settings.Settings, err error) {
	settings.VPN, err = readVPN()
	if err != nil {
		return settings, err
	}

	settings.HTTPProxy, err = readHTTPProxy()
	if err != nil {
		return settings, err
	}

	return settings, nil
}
