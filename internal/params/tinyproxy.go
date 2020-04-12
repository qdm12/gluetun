package params

import (
	"strconv"

	libparams "github.com/qdm12/golibs/params"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

// GetTinyProxy obtains if TinyProxy is on from the environment variable
// TINYPROXY, and using PROXY as a retro-compatibility name
func (p *reader) GetTinyProxy() (activated bool, err error) {
	// Retro-compatibility
	s, err := p.envParams.GetEnv("PROXY")
	if err != nil {
		return false, err
	} else if len(s) != 0 {
		p.logger.Warn("You are using the old environment variable PROXY, please consider changing it to TINYPROXY")
		return p.envParams.GetOnOff("PROXY", libparams.Compulsory())
	}
	return p.envParams.GetOnOff("TINYPROXY", libparams.Default("off"))
}

// GetTinyProxyLog obtains the TinyProxy log level from the environment variable
// TINYPROXY_LOG, and using PROXY_LOG_LEVEL as a retro-compatibility name
func (p *reader) GetTinyProxyLog() (models.TinyProxyLogLevel, error) {
	// Retro-compatibility
	s, err := p.envParams.GetEnv("PROXY_LOG_LEVEL")
	if err != nil {
		return models.TinyProxyLogLevel(s), err
	} else if len(s) != 0 {
		p.logger.Warn("You are using the old environment variable PROXY_LOG_LEVEL, please consider changing it to TINYPROXY_LOG")
		s, err = p.envParams.GetValueIfInside("PROXY_LOG_LEVEL", []string{"Info", "Connect", "Notice", "Warning", "Error", "Critical"}, libparams.Compulsory())
		return models.TinyProxyLogLevel(s), err
	}
	s, err = p.envParams.GetValueIfInside("TINYPROXY_LOG", []string{"Info", "Connect", "Notice", "Warning", "Error", "Critical"}, libparams.Default("Connect"))
	return models.TinyProxyLogLevel(s), err
}

// GetTinyProxyPort obtains the TinyProxy listening port from the environment variable
// TINYPROXY_PORT, and using PROXY_PORT as a retro-compatibility name
func (p *reader) GetTinyProxyPort() (port uint16, err error) {
	// Retro-compatibility
	portStr, err := p.envParams.GetEnv("PROXY_PORT")
	switch {
	case err != nil:
		return 0, err
	case len(portStr) != 0:
		p.logger.Warn("You are using the old environment variable PROXY_PORT, please consider changing it to TINYPROXY_PORT")
	default:
		portStr, err = p.envParams.GetEnv("TINYPROXY_PORT", libparams.Default("8888"))
		if err != nil {
			return 0, err
		}
	}
	if err := p.verifier.VerifyPort(portStr); err != nil {
		return 0, err
	}
	portUint64, err := strconv.ParseUint(portStr, 10, 16)
	return uint16(portUint64), err
}

// GetTinyProxyUser obtains the TinyProxy server user from the environment variable
// TINYPROXY_USER, and using PROXY_USER as a retro-compatibility name
func (p *reader) GetTinyProxyUser() (user string, err error) {
	defer func() {
		unsetErr := p.unsetEnv("PROXY_USER")
		if err == nil {
			err = unsetErr
		}
	}()
	defer func() {
		unsetErr := p.unsetEnv("TINYPROXY_USER")
		if err == nil {
			err = unsetErr
		}
	}()
	// Retro-compatibility
	user, err = p.envParams.GetEnv("PROXY_USER", libparams.CaseSensitiveValue())
	if err != nil {
		return user, err
	}
	if len(user) != 0 {
		p.logger.Warn("You are using the old environment variable PROXY_USER, please consider changing it to TINYPROXY_USER")
		return user, nil
	}
	return p.envParams.GetEnv("TINYPROXY_USER", libparams.CaseSensitiveValue())
}

// GetTinyProxyPassword obtains the TinyProxy server password from the environment variable
// TINYPROXY_PASSWORD, and using PROXY_PASSWORD as a retro-compatibility name
func (p *reader) GetTinyProxyPassword() (password string, err error) {
	defer func() {
		unsetErr := p.unsetEnv("PROXY_PASSWORD")
		if err == nil {
			err = unsetErr
		}
	}()
	defer func() {
		unsetErr := p.unsetEnv("TINYPROXY_PASSWORD")
		if err == nil {
			err = unsetErr
		}
	}()

	// Retro-compatibility
	password, err = p.envParams.GetEnv("PROXY_PASSWORD", libparams.CaseSensitiveValue())
	if err != nil {
		return password, err
	}
	if len(password) != 0 {
		p.logger.Warn("You are using the old environment variable PROXY_PASSWORD, please consider changing it to TINYPROXY_PASSWORD")
		return password, nil
	}
	return p.envParams.GetEnv("TINYPROXY_PASSWORD", libparams.CaseSensitiveValue())
}
