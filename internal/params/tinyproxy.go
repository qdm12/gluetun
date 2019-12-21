package params

import (
	"fmt"
	"github.com/qdm12/golibs/logging"
	libparams "github.com/qdm12/golibs/params"
	"github.com/qdm12/golibs/verification"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
)

// GetTinyProxy obtains if TinyProxy is on from the environment variable
// TINYPROXY, and using PROXY as a retro-compatibility name
func GetTinyProxy() (activated bool, err error) {
	// Retro-compatibility
	s := libparams.GetEnv("PROXY", "")
	if len(s) != 0 {
		logging.Warn("You are using the old environment variable PROXY, please consider changing it to TINYPROXY")
		if s == "on" {
			return true, nil
		} else if s == "off" {
			return false, nil
		}
		return false, fmt.Errorf("Environment variable PROXY can only be \"on\" or \"off\"")
	}
	return libparams.GetOnOff("TINYPROXY", false)
}

// GetTinyProxyLog obtains the TinyProxy log level from the environment variable
// TINYPROXY_LOG, and using PROXY_LOG_LEVEL as a retro-compatibility name
func GetTinyProxyLog() (constants.TinyProxyLogLevel, error) {
	// Retro-compatibility
	if libparams.GetEnv("PROXY_LOG_LEVEL", "") != "" {
		logging.Warn("You are using the old environment variable PROXY_LOG_LEVEL, please consider changing it to TINYPROXY_LOG")
		s, err := libparams.GetValueIfInside("PROXY_LOG_LEVEL", []string{"info", "warning", "error", "critical"}, true, "")
		return constants.TinyProxyLogLevel(s), err
	}
	s, err := libparams.GetValueIfInside("TINYPROXY_LOG", []string{"info", "warning", "error", "critical"}, false, "info")
	return constants.TinyProxyLogLevel(s), err
}

// GetTinyProxyPort obtains the TinyProxy listening port from the environment variable
// TINYPROXY_PORT, and using PROXY_PORT as a retro-compatibility name
func GetTinyProxyPort() (port string, err error) {
	// Retro-compatibility
	port = libparams.GetEnv("PROXY_PORT", "")
	if len(port) != 0 {
		logging.Warn("You are using the old environment variable PROXY_PORT, please consider changing it to TINYPROXY_PORT")
	} else {
		port = libparams.GetEnv("TINYPROXY_PORT", "")
	}
	return port, verification.VerifyPort(port)
}

// GetTinyProxyUser obtains the TinyProxy server user from the environment variable
// TINYPROXY_USER, and using PROXY_USER as a retro-compatibility name
func GetTinyProxyUser() (user string) {
	// Retro-compatibility
	user = libparams.GetEnv("PROXY_USER", "")
	if len(user) != 0 {
		logging.Warn("You are using the old environment variable PROXY_USER, please consider changing it to TINYPROXY_USER")
		return user
	}
	return libparams.GetEnv("TINYPROXY_USER", "")
}

// GetTinyProxyPassword obtains the TinyProxy server password from the environment variable
// TINYPROXY_PASSWORD, and using PROXY_PASSWORD as a retro-compatibility name
func GetTinyProxyPassword() (password string) {
	// Retro-compatibility
	password = libparams.GetEnv("PROXY_PASSWORD", "")
	if len(password) != 0 {
		logging.Warn("You are using the old environment variable PROXY_PASSWORD, please consider changing it to TINYPROXY_PASSWORD")
		return password
	}
	return libparams.GetEnv("TINYPROXY_PASSWORD", "")
}
