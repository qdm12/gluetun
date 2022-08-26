package secrets

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
)

type Source struct{}

func New() *Source {
	return &Source{}
}

func (s *Source) String() string { return "secret files" }

func (s *Source) Read() (settings settings.Settings, err error) {
	settings.VPN, err = readVPN()
	if err != nil {
		return settings, err
	}

	settings.HTTPProxy, err = readHTTPProxy()
	if err != nil {
		return settings, err
	}

	settings.Shadowsocks, err = readShadowsocks()
	if err != nil {
		return settings, err
	}

	return settings, nil
}
