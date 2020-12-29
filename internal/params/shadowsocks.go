package params

import (
	"strconv"

	libparams "github.com/qdm12/golibs/params"
)

// GetShadowSocks obtains if ShadowSocks is on from the environment variable
// SHADOWSOCKS.
func (r *reader) GetShadowSocks() (activated bool, err error) {
	return r.envParams.GetOnOff("SHADOWSOCKS", libparams.Default("off"))
}

// GetShadowSocksLog obtains the ShadowSocks log level from the environment variable
// SHADOWSOCKS_LOG.
func (r *reader) GetShadowSocksLog() (activated bool, err error) {
	return r.envParams.GetOnOff("SHADOWSOCKS_LOG", libparams.Default("off"))
}

// GetShadowSocksPort obtains the ShadowSocks listening port from the environment variable
// SHADOWSOCKS_PORT.
func (r *reader) GetShadowSocksPort() (port uint16, err error) {
	portStr, err := r.envParams.GetEnv("SHADOWSOCKS_PORT", libparams.Default("8388"))
	if err != nil {
		return 0, err
	}
	if err := r.verifier.VerifyPort(portStr); err != nil {
		return 0, err
	}
	portUint64, err := strconv.ParseUint(portStr, 10, 16)
	return uint16(portUint64), err
}

// GetShadowSocksPassword obtains the ShadowSocks server password.
// It first tries to use the SHADOWSOCKS_PASSWORD environment variable (easier for the end user)
// and then tries to read from the secret file shadowsocks_password if nothing was found.
func (r *reader) GetShadowSocksPassword() (password string, err error) {
	const compulsory = false
	return r.getFromEnvOrSecretFile("SHADOWSOCKS_PASSWORD", compulsory, nil)
}

// GetShadowSocksMethod obtains the ShadowSocks method to use from the environment variable
// SHADOWSOCKS_METHOD.
func (r *reader) GetShadowSocksMethod() (method string, err error) {
	return r.envParams.GetEnv("SHADOWSOCKS_METHOD", libparams.Default("chacha20-ietf-poly1305"))
}
