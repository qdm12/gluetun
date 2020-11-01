package params

import (
	"strings"

	libparams "github.com/qdm12/golibs/params"
)

// GetHTTPProxy obtains if the HTTP proxy is on from the environment variable
// HTTPPROXY, and using PROXY and TINYPROXY as retro-compatibility names.
func (r *reader) GetHTTPProxy() (enabled bool, err error) {
	retroKeysOption := libparams.RetroKeys(
		[]string{"TINYPROXY", "PROXY"},
		r.onRetroActive,
	)
	return r.envParams.GetOnOff("HTTPPROXY", retroKeysOption, libparams.Default("off"))
}

// GetHTTPProxyLog obtains the if http proxy requests should be logged from
// the environment variable HTTPPROXY_LOG, and using PROXY_LOG_LEVEL and
// TINYPROXY_LOG as retro-compatibility names.
func (r *reader) GetHTTPProxyLog() (log bool, err error) {
	s, _ := r.envParams.GetEnv("HTTPPROXY_LOG")
	if len(s) == 0 {
		s, _ = r.envParams.GetEnv("PROXY_LOG_LEVEL")
		if len(s) == 0 {
			s, _ = r.envParams.GetEnv("TINYPROXY_LOG")
			if len(s) == 0 {
				return false, nil // default log disabled
			}
		}
		switch strings.ToLower(s) {
		case "info", "connect", "notice":
			return true, nil
		default:
			return false, nil
		}
	}
	return r.envParams.GetOnOff("HTTPPROXY_LOG", libparams.Default("off"))
}

// GetHTTPProxyPort obtains the HTTP proxy listening port from the environment variable
// HTTPPROXY_PORT, and using PROXY_PORT and TINYPROXY_PORT as retro-compatibility names.
func (r *reader) GetHTTPProxyPort() (port uint16, err error) {
	retroKeysOption := libparams.RetroKeys(
		[]string{"TINYPROXY_PORT", "PROXY_PORT"},
		r.onRetroActive,
	)
	return r.envParams.GetPort("HTTPPROXY_PORT", retroKeysOption, libparams.Default("8888"))
}

// GetHTTPProxyUser obtains the HTTP proxy server user from the environment variable
// HTTPPROXY_USER, and using TINYPROXY_USER and PROXY_USER as retro-compatibility names.
func (r *reader) GetHTTPProxyUser() (user string, err error) {
	retroKeysOption := libparams.RetroKeys(
		[]string{"TINYPROXY_USER", "PROXY_USER"},
		r.onRetroActive,
	)
	return r.envParams.GetEnv("HTTPPROXY_USER",
		retroKeysOption, libparams.CaseSensitiveValue(), libparams.Unset())
}

// GetHTTPProxyPassword obtains the HTTP proxy server password from the environment variable
// HTTPPROXY_PASSWORD, and using TINYPROXY_PASSWORD and PROXY_PASSWORD as retro-compatibility names.
func (r *reader) GetHTTPProxyPassword() (password string, err error) {
	retroKeysOption := libparams.RetroKeys(
		[]string{"TINYPROXY_PASSWORD", "PROXY_PASSWORD"},
		r.onRetroActive,
	)
	return r.envParams.GetEnv("HTTPPROXY_PASSWORD",
		retroKeysOption, libparams.CaseSensitiveValue(), libparams.Unset())
}

// GetHTTPProxyStealth obtains the HTTP proxy server stealth mode
// from the environment variable HTTPPROXY_STEALTH.
func (r *reader) GetHTTPProxyStealth() (stealth bool, err error) {
	return r.envParams.GetOnOff("HTTPPROXY_STEALTH", libparams.Default("off"))
}
