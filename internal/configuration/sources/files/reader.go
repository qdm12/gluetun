package files

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
)

type Source struct {
	wireguardConfigPath string
}

func New() *Source {
	const wireguardConfigPath = "/gluetun/wireguard/wg0.conf"
	return &Source{
		wireguardConfigPath: wireguardConfigPath,
	}
}

func (s *Source) String() string { return "files" }

func (s *Source) Read() (settings settings.Settings, err error) {
	settings.VPN, err = s.readVPN()
	if err != nil {
		return settings, err
	}

	settings.System, err = s.readSystem()
	if err != nil {
		return settings, err
	}

	return settings, nil
}
