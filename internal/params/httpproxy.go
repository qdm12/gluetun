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
	return r.env.OnOff("HTTPPROXY", retroKeysOption, libparams.Default("off"))
}

// GetHTTPProxyLog obtains the if http proxy requests should be logged from
// the environment variable HTTPPROXY_LOG, and using PROXY_LOG_LEVEL and
// TINYPROXY_LOG as retro-compatibility names.
func (r *reader) GetHTTPProxyLog() (log bool, err error) {
	s, _ := r.env.Get("HTTPPROXY_LOG")
	if len(s) == 0 {
		s, _ = r.env.Get("PROXY_LOG_LEVEL")
		if len(s) == 0 {
			s, _ = r.env.Get("TINYPROXY_LOG")
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
	return r.env.OnOff("HTTPPROXY_LOG", libparams.Default("off"))
}

// GetHTTPProxyPort obtains the HTTP proxy listening port from the environment variable
// HTTPPROXY_PORT, and using PROXY_PORT and TINYPROXY_PORT as retro-compatibility names.
func (r *reader) GetHTTPProxyPort() (port uint16, warning string, err error) {
	retroKeysOption := libparams.RetroKeys(
		[]string{"TINYPROXY_PORT", "PROXY_PORT"},
		r.onRetroActive,
	)
	return r.env.ListeningPort("HTTPPROXY_PORT", retroKeysOption, libparams.Default("8888"))
}

// GetHTTPProxyUser obtains the HTTP proxy server user.
// It first tries to use the HTTPPROXY_USER environment variable (easier for the end user)
// and then tries to read from the secret file httpproxy_user if nothing was found.
func (r *reader) GetHTTPProxyUser() (user string, err error) {
	const compulsory = false
	return r.getFromEnvOrSecretFile(
		"HTTPPROXY_USER",
		compulsory,
		[]string{"TINYPROXY_USER", "PROXY_USER"},
	)
}

// GetHTTPProxyPassword obtains the HTTP proxy server password.
// It first tries to use the HTTPPROXY_PASSWORD environment variable (easier for the end user)
// and then tries to read from the secret file httpproxy_password if nothing was found.
func (r *reader) GetHTTPProxyPassword() (password string, err error) {
	const compulsory = false
	return r.getFromEnvOrSecretFile(
		"HTTPPROXY_USER",
		compulsory,
		[]string{"TINYPROXY_PASSWORD", "PROXY_PASSWORD"},
	)
}

// GetHTTPProxyStealth obtains the HTTP proxy server stealth mode
// from the environment variable HTTPPROXY_STEALTH.
func (r *reader) GetHTTPProxyStealth() (stealth bool, err error) {
	return r.env.OnOff("HTTPPROXY_STEALTH", libparams.Default("off"))
}
