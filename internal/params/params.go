package params

import (
	"fmt"
	libparams "github.com/qdm12/golibs/params"
	"github.com/qdm12/golibs/verification"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
)

// GetPortForwarding obtains if port forwarding on the VPN provider server
// side is enabled or not from the environment variable PORT_FORWARDING
func GetPortForwarding() (activated bool, err error) {
	s := libparams.GetEnv("PORT_FORWARDING", "off")
	if s == "false" || s == "off" {
		return false, nil
	} else if s == "true" || s == "on" {
		return true, nil
	}
	return false, fmt.Errorf("PORT_FORWARDING can only be \"on\" or \"off\"")
}

// GetTinyProxy obtains if TinyProxy is on from the environment variable
// TINYPROXY, and using PROXY as a retro-compatibility name
func GetTinyProxy() (activated bool, err error) {
	// Retro-compatibility
	s := libparams.GetEnv("PROXY", "")

	return libparams.GetOnOff("TINYPROXY", s == "on")
}

// GetTinyProxyLog obtains the TinyProxy log level from the environment variable
// TINYPROXY_LOG, and using PROXY_LOG_LEVEL as a retro-compatibility name
func GetTinyProxyLog() (logLevel constants.TinyProxyLogLevel, err error) {
	s := libparams.GetEnv("PROXY_LOG_LEVEL", "") // Retro-compatibility
	s = libparams.GetEnv("TINYPROXY_LOG", s)
	return constants.ParseTinyProxyLogLevel(s)
}

// GetTinyProxyPort obtains the TinyProxy listening port from the environment variable
// TINYPROXY_PORT, and using PROXY_PORT as a retro-compatibility name
func GetTinyProxyPort() (port string, err error) {
	s := libparams.GetEnv("PROXY_PORT", "") // Retro-compatibility
	s = libparams.GetEnv("TINYPROXY_PORT", s)
	return s, verification.VerifyPort(s)
}
