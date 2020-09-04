package settings

import (
	"fmt"
	"strings"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/params"
)

// TinyProxy contains settings to configure TinyProxy
type TinyProxy struct {
	User     string
	Password string
	LogLevel models.TinyProxyLogLevel
	Port     uint16
	Enabled  bool
}

func (t *TinyProxy) String() string {
	if !t.Enabled {
		return "TinyProxy settings: disabled"
	}
	auth := disabled
	if t.User != "" {
		auth = enabled
	}
	settingsList := []string{
		fmt.Sprintf("Port: %d", t.Port),
		"Authentication: " + auth,
		"Log level: " + string(t.LogLevel),
	}
	return "TinyProxy settings:\n" + strings.Join(settingsList, "\n |--")
}

// GetTinyProxySettings obtains TinyProxy settings from environment variables using the params package.
func GetTinyProxySettings(paramsReader params.Reader) (settings TinyProxy, err error) {
	settings.Enabled, err = paramsReader.GetTinyProxy()
	if err != nil || !settings.Enabled {
		return settings, err
	}
	settings.User, err = paramsReader.GetTinyProxyUser()
	if err != nil {
		return settings, err
	}
	settings.Password, err = paramsReader.GetTinyProxyPassword()
	if err != nil {
		return settings, err
	}
	settings.Port, err = paramsReader.GetTinyProxyPort()
	if err != nil {
		return settings, err
	}
	settings.LogLevel, err = paramsReader.GetTinyProxyLog()
	if err != nil {
		return settings, err
	}
	return settings, nil
}
