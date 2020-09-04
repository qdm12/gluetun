package settings

import (
	"fmt"
	"strings"

	"github.com/qdm12/gluetun/internal/params"
)

// ShadowSocks contains settings to configure the Shadowsocks server
type ShadowSocks struct {
	Method   string
	Password string
	Port     uint16
	Enabled  bool
	Log      bool
}

func (s *ShadowSocks) String() string {
	if !s.Enabled {
		return "ShadowSocks settings: disabled"
	}
	log := disabled
	if s.Log {
		log = enabled
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
func GetShadowSocksSettings(paramsReader params.Reader) (settings ShadowSocks, err error) {
	settings.Enabled, err = paramsReader.GetShadowSocks()
	if err != nil || !settings.Enabled {
		return settings, err
	}
	settings.Port, err = paramsReader.GetShadowSocksPort()
	if err != nil {
		return settings, err
	}
	settings.Password, err = paramsReader.GetShadowSocksPassword()
	if err != nil {
		return settings, err
	}
	settings.Log, err = paramsReader.GetShadowSocksLog()
	if err != nil {
		return settings, err
	}
	settings.Method, err = paramsReader.GetShadowSocksMethod()
	if err != nil {
		return settings, err
	}
	return settings, nil
}
