package settings

import (
	"fmt"
	"strings"

	"github.com/qdm12/gluetun/internal/params"
)

// HTTPProxy contains settings to configure the HTTP proxy.
type HTTPProxy struct { //nolint:maligned
	Enabled  bool
	Port     uint16
	User     string
	Password string
	Stealth  bool
	Log      bool
}

func (h *HTTPProxy) String() string {
	if !h.Enabled {
		return "HTTPProxy settings: disabled"
	}
	auth, log, stealth := disabled, disabled, disabled
	if h.User != "" {
		auth = enabled
	}
	if h.Log {
		log = enabled
	}
	if h.Stealth {
		stealth = enabled
	}
	settingsList := []string{
		fmt.Sprintf("Port: %d", h.Port),
		"Authentication: " + auth,
		"Stealth: " + stealth,
		"Log: " + log,
	}
	return "HTTP proxy settings:\n" + strings.Join(settingsList, "\n |--")
}

// GetHTTPProxySettings obtains HTTPProxy settings from environment variables using the params package.
func GetHTTPProxySettings(paramsReader params.Reader) (settings HTTPProxy, err error) {
	settings.Enabled, err = paramsReader.GetHTTPProxy()
	if err != nil || !settings.Enabled {
		return settings, err
	}
	settings.Port, err = paramsReader.GetHTTPProxyPort()
	if err != nil {
		return settings, err
	}
	settings.User, err = paramsReader.GetHTTPProxyUser()
	if err != nil {
		return settings, err
	}
	settings.Password, err = paramsReader.GetHTTPProxyPassword()
	if err != nil {
		return settings, err
	}
	settings.Stealth, err = paramsReader.GetHTTPProxyStealth()
	if err != nil {
		return settings, err
	}
	settings.Log, err = paramsReader.GetHTTPProxyLog()
	if err != nil {
		return settings, err
	}
	return settings, nil
}
