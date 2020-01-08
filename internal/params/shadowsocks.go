package params

import (
	libparams "github.com/qdm12/golibs/params"
	"github.com/qdm12/golibs/verification"
)

// GetShadowSocks obtains if ShadowSocks is on from the environment variable
// SHADOWSOCKS
func GetShadowSocks(envParams libparams.EnvParams) (activated bool, err error) {
	return envParams.GetOnOff("SHADOWSOCKS", libparams.Default("off"))
}

// GetShadowSocksLog obtains the ShadowSocks log level from the environment variable
// TINYPROXY_LOG, and using PROXY_LOG_LEVEL as a retro-compatibility name
func GetShadowSocksLog(envParams libparams.EnvParams) (activated bool, err error) {
	return envParams.GetOnOff("SHADOWSOCKS_LOG", libparams.Default("off"))
}

// GetShadowSocksPort obtains the ShadowSocks listening port from the environment variable
// SHADOWSOCKS_PORT
func GetShadowSocksPort(envParams libparams.EnvParams) (port string, err error) {
	port, err = envParams.GetEnv("SHADOWSOCKS_PORT", libparams.Default("8388"))
	if err != nil {
		return port, err
	}
	return port, verification.VerifyPort(port)
}

// GetShadowSocksPassword obtains the ShadowSocks server password from the environment variable
// SHADOWSOCKS_PASSWORD
func GetShadowSocksPassword(envParams libparams.EnvParams) (password string, err error) {
	return envParams.GetEnv("SHADOWSOCKS_PASSWORD")
}
