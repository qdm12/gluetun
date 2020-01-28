package params

import (
	libparams "github.com/qdm12/golibs/params"
)

// GetShadowSocks obtains if ShadowSocks is on from the environment variable
// SHADOWSOCKS
func (p *paramsReader) GetShadowSocks() (activated bool, err error) {
	return p.envParams.GetOnOff("SHADOWSOCKS", libparams.Default("off"))
}

// GetShadowSocksLog obtains the ShadowSocks log level from the environment variable
// TINYPROXY_LOG, and using PROXY_LOG_LEVEL as a retro-compatibility name
func (p *paramsReader) GetShadowSocksLog() (activated bool, err error) {
	return p.envParams.GetOnOff("SHADOWSOCKS_LOG", libparams.Default("off"))
}

// GetShadowSocksPort obtains the ShadowSocks listening port from the environment variable
// SHADOWSOCKS_PORT
func (p *paramsReader) GetShadowSocksPort() (port string, err error) {
	port, err = p.envParams.GetEnv("SHADOWSOCKS_PORT", libparams.Default("8388"))
	if err != nil {
		return port, err
	}
	return port, p.verifier.VerifyPort(port)
}

// GetShadowSocksPassword obtains the ShadowSocks server password from the environment variable
// SHADOWSOCKS_PASSWORD
func (p *paramsReader) GetShadowSocksPassword() (password string, err error) {
	return p.envParams.GetEnv("SHADOWSOCKS_PASSWORD")
}
