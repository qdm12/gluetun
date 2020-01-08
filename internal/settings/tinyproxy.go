package settings

import (
	"strings"

	libparams "github.com/qdm12/golibs/params"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/params"
)

// TinyProxy contains settings to configure TinyProxy
type TinyProxy struct {
	Enabled  bool
	User     string
	Password string
	Port     string
	LogLevel constants.TinyProxyLogLevel
}

func (t *TinyProxy) String() string {
	if !t.Enabled {
		return "TinyProxy settings: disabled"
	}
	auth := "disabled"
	if t.User != "" {
		auth = "enabled"
	}
	settingsList := []string{
		"Port: " + t.Port,
		"Authentication: " + auth,
		"Log level: " + string(t.LogLevel),
	}
	return "TinyProxy settings:\n" + strings.Join(settingsList, "\n |--")
}

// GetTinyProxySettings obtains TinyProxy settings from environment variables using the params package.
func GetTinyProxySettings(envParams libparams.EnvParams) (settings TinyProxy, err error) {
	settings.Enabled, err = params.GetTinyProxy(envParams)
	if err != nil || !settings.Enabled {
		return settings, err
	}
	settings.User, err = params.GetTinyProxyUser(envParams)
	if err != nil {
		return settings, err
	}
	settings.Password, err = params.GetTinyProxyPassword(envParams)
	if err != nil {
		return settings, err
	}
	settings.Port, err = params.GetTinyProxyPort(envParams)
	if err != nil {
		return settings, err
	}
	settings.LogLevel, err = params.GetTinyProxyLog(envParams)
	if err != nil {
		return settings, err
	}
	return settings, nil
}
