package params

import (
	libparams "github.com/qdm12/golibs/params"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
)

// GetTinyProxy obtains if TinyProxy is on from the environment variable
// TINYPROXY, and using PROXY as a retro-compatibility name
func (p *paramsReader) GetTinyProxy() (activated bool, err error) {
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
func (p *paramsReader) GetTinyProxyLog() (constants.TinyProxyLogLevel, error) {
	// Retro-compatibility
	s, err := p.envParams.GetEnv("PROXY_LOG_LEVEL")
	if err != nil {
		return constants.TinyProxyLogLevel(s), err
	} else if len(s) != 0 {
		p.logger.Warn("You are using the old environment variable PROXY_LOG_LEVEL, please consider changing it to TINYPROXY_LOG")
		s, err = p.envParams.GetValueIfInside("PROXY_LOG_LEVEL", []string{"info", "warning", "error", "critical"}, libparams.Compulsory())
		return constants.TinyProxyLogLevel(s), err
	}
	s, err = p.envParams.GetValueIfInside("TINYPROXY_LOG", []string{"info", "warning", "error", "critical"}, libparams.Default("info"))
	return constants.TinyProxyLogLevel(s), err
}

// GetTinyProxyPort obtains the TinyProxy listening port from the environment variable
// TINYPROXY_PORT, and using PROXY_PORT as a retro-compatibility name
func (p *paramsReader) GetTinyProxyPort() (port string, err error) {
	// Retro-compatibility
	port, err = p.envParams.GetEnv("PROXY_PORT")
	if err != nil {
		return port, err
	} else if len(port) != 0 {
		p.logger.Warn("You are using the old environment variable PROXY_PORT, please consider changing it to TINYPROXY_PORT")
	} else {
		port, err = p.envParams.GetEnv("TINYPROXY_PORT", libparams.Default("8888"))
		if err != nil {
			return port, err
		}
	}
	return port, p.verifier.VerifyPort(port)
}

// GetTinyProxyUser obtains the TinyProxy server user from the environment variable
// TINYPROXY_USER, and using PROXY_USER as a retro-compatibility name
func (p *paramsReader) GetTinyProxyUser() (user string, err error) {
	// Retro-compatibility
	user, err = p.envParams.GetEnv("PROXY_USER")
	if err != nil {
		return user, err
	}
	if len(user) != 0 {
		p.logger.Warn("You are using the old environment variable PROXY_USER, please consider changing it to TINYPROXY_USER")
		return user, nil
	}
	return p.envParams.GetEnv("TINYPROXY_USER")
}

// GetTinyProxyPassword obtains the TinyProxy server password from the environment variable
// TINYPROXY_PASSWORD, and using PROXY_PASSWORD as a retro-compatibility name
func (p *paramsReader) GetTinyProxyPassword() (password string, err error) {
	// Retro-compatibility
	password, err = p.envParams.GetEnv("PROXY_PASSWORD")
	if err != nil {
		return password, err
	}
	if len(password) != 0 {
		p.logger.Warn("You are using the old environment variable PROXY_PASSWORD, please consider changing it to TINYPROXY_PASSWORD")
		return password, nil
	}
	return p.envParams.GetEnv("TINYPROXY_PASSWORD")
}
