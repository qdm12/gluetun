package params

import (
	libparams "github.com/qdm12/golibs/params"
	"github.com/qdm12/golibs/verification"
)

// GetShadowSocks obtains if ShadowSocks is on from the environment variable
// SHADOWSOCKS
func GetShadowSocks() (activated bool, err error) {
	return libparams.GetOnOff("SHADOWSOCKS", false)
}

// GetShadowSocksLog obtains the ShadowSocks log level from the environment variable
// TINYPROXY_LOG, and using PROXY_LOG_LEVEL as a retro-compatibility name
func GetShadowSocksLog() (activated bool, err error) {
	return libparams.GetOnOff("SHADOWSOCKS_LOG", false)
}

// GetShadowSocksPort obtains the ShadowSocks listening port from the environment variable
// SHADOWSOCKS_PORT
func GetShadowSocksPort() (port string, err error) {
	port = libparams.GetEnv("SHADOWSOCKS_PORT", "")
	return port, verification.VerifyPort(port)
}

// GetShadowSocksPassword obtains the ShadowSocks server password from the environment variable
// SHADOWSOCKS_PASSWORD
func GetShadowSocksPassword() (password string) {
	return libparams.GetEnv("SHADOWSOCKS_PASSWORD", "")
}
