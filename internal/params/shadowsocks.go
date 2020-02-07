package params

import (
	"strconv"

	libparams "github.com/qdm12/golibs/params"
)

// GetShadowSocks obtains if ShadowSocks is on from the environment variable
// SHADOWSOCKS
func (p *paramsReader) GetShadowSocks() (activated bool, err error) {
	return p.envParams.GetOnOff("SHADOWSOCKS", libparams.Default("off"))
}

// GetShadowSocksLog obtains the ShadowSocks log level from the environment variable
// SHADOWSOCKS_LOG
func (p *paramsReader) GetShadowSocksLog() (activated bool, err error) {
	return p.envParams.GetOnOff("SHADOWSOCKS_LOG", libparams.Default("off"))
}

// GetShadowSocksPort obtains the ShadowSocks listening port from the environment variable
// SHADOWSOCKS_PORT
func (p *paramsReader) GetShadowSocksPort() (port uint16, err error) {
	portStr, err := p.envParams.GetEnv("SHADOWSOCKS_PORT", libparams.Default("8388"))
	if err != nil {
		return 0, err
	}
	if err := p.verifier.VerifyPort(portStr); err != nil {
		return 0, err
	}
	portUint64, err := strconv.ParseUint(portStr, 10, 16)
	return uint16(portUint64), err
}

// GetShadowSocksPassword obtains the ShadowSocks server password from the environment variable
// SHADOWSOCKS_PASSWORD
func (p *paramsReader) GetShadowSocksPassword() (password string, err error) {
	defer p.unsetEnv("SHADOWSOCKS_PASSWORD")
	return p.envParams.GetEnv("SHADOWSOCKS_PASSWORD")
}
