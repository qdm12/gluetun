package settings

import (
	"strings"

	libparams "github.com/qdm12/golibs/params"
	"github.com/qdm12/private-internet-access-docker/internal/params"
)

// ShadowSocks contains settings to configure the Shadowsocks server
type ShadowSocks struct {
	Enabled  bool
	Password string
	Log      bool
	Port     string
}

func (s *ShadowSocks) String() string {
	if !s.Enabled {
		return "ShadowSocks settings: disabled"
	}
	settingsList := []string{
		"Port: " + s.Port,
	}
	return "ShadowSocks settings:\n" + strings.Join(settingsList, "\n |--")
}

// GetShadowSocksSettings obtains ShadowSocks settings from environment variables using the params package.
func GetShadowSocksSettings(envParams libparams.EnvParams) (settings ShadowSocks, err error) {
	settings.Enabled, err = params.GetShadowSocks(envParams)
	if err != nil || !settings.Enabled {
		return settings, err
	}
	settings.Port, err = params.GetShadowSocksPort(envParams)
	if err != nil {
		return settings, err
	}
	settings.Password, err = params.GetShadowSocksPassword(envParams)
	if err != nil {
		return settings, err
	}
	settings.Log, err = params.GetShadowSocksLog(envParams)
	if err != nil {
		return settings, err
	}
	return settings, nil
}
