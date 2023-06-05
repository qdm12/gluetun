package secrets

import (
	"os"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gosettings/sources/env"
)

type Source struct {
	env env.Env
}

func New() *Source {
	handleDeprecatedKey := (func(deprecatedKey, newKey string))(nil)
	return &Source{
		env: *env.New(os.Environ(), handleDeprecatedKey),
	}
}

func (s *Source) String() string { return "secret files" }

func (s *Source) Read() (settings settings.Settings, err error) {
	settings.VPN, err = s.readVPN()
	if err != nil {
		return settings, err
	}

	settings.HTTPProxy, err = s.readHTTPProxy()
	if err != nil {
		return settings, err
	}

	settings.Shadowsocks, err = s.readShadowsocks()
	if err != nil {
		return settings, err
	}

	return settings, nil
}
