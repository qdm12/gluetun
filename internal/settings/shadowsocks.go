package settings

import (
	"fmt"
	"strings"

	"github.com/qdm12/private-internet-access-docker/internal/params"
)

// ShadowSocks contains settings to configure the Shadowsocks server
type ShadowSocks struct {
	Enabled  bool
	Password string
	Log      bool
	Port     uint16
	Method   string
}

func (s *ShadowSocks) String() string {
	if !s.Enabled {
		return "ShadowSocks settings: disabled"
	}
	log := "disabled"
	if s.Log {
		log = "enabled"
	}
	settingsList := []string{
		"ShadowSocks settings:",
		"Password: [redacted]",
		"Log: " + log,
		fmt.Sprintf("Port: %d", s.Port),
		"Method: " + s.Method,
	}
	return strings.Join(settingsList, "\n |--")
}

// GetShadowSocksSettings obtains ShadowSocks settings from environment variables using the params package.
func GetShadowSocksSettings(params params.ParamsReader) (settings ShadowSocks, err error) {
	settings.Enabled, err = params.GetShadowSocks()
	if err != nil || !settings.Enabled {
		return settings, err
	}
	settings.Port, err = params.GetShadowSocksPort()
	if err != nil {
		return settings, err
	}
	settings.Password, err = params.GetShadowSocksPassword()
	if err != nil {
		return settings, err
	}
	settings.Log, err = params.GetShadowSocksLog()
	if err != nil {
		return settings, err
	}
	settings.Method, err = params.GetShadowSocksMethod()
	if err != nil {
		return settings, err
	}
	return settings, nil
}
