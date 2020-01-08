package params

import (
	"github.com/qdm12/golibs/logging"
	libparams "github.com/qdm12/golibs/params"
	"github.com/qdm12/golibs/verification"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
)

// GetTinyProxy obtains if TinyProxy is on from the environment variable
// TINYPROXY, and using PROXY as a retro-compatibility name
func GetTinyProxy(envParams libparams.EnvParams) (activated bool, err error) {
	// Retro-compatibility
	s, err := envParams.GetEnv("PROXY")
	if err != nil {
		return false, err
	} else if len(s) != 0 {
		logging.Warn("You are using the old environment variable PROXY, please consider changing it to TINYPROXY")
		return envParams.GetOnOff("PROXY", libparams.Compulsory())
	}
	return envParams.GetOnOff("TINYPROXY", libparams.Default("off"))
}

// GetTinyProxyLog obtains the TinyProxy log level from the environment variable
// TINYPROXY_LOG, and using PROXY_LOG_LEVEL as a retro-compatibility name
func GetTinyProxyLog(envParams libparams.EnvParams) (constants.TinyProxyLogLevel, error) {
	// Retro-compatibility
	s, err := envParams.GetEnv("PROXY_LOG_LEVEL")
	if err != nil {
		return constants.TinyProxyLogLevel(s), err
	} else if len(s) != 0 {
		logging.Warn("You are using the old environment variable PROXY_LOG_LEVEL, please consider changing it to TINYPROXY_LOG")
		s, err = envParams.GetValueIfInside("PROXY_LOG_LEVEL", []string{"info", "warning", "error", "critical"}, libparams.Compulsory())
		return constants.TinyProxyLogLevel(s), err
	}
	s, err = envParams.GetValueIfInside("TINYPROXY_LOG", []string{"info", "warning", "error", "critical"}, libparams.Default("info"))
	return constants.TinyProxyLogLevel(s), err
}

// GetTinyProxyPort obtains the TinyProxy listening port from the environment variable
// TINYPROXY_PORT, and using PROXY_PORT as a retro-compatibility name
func GetTinyProxyPort(envParams libparams.EnvParams) (port string, err error) {
	// Retro-compatibility
	port, err = envParams.GetEnv("PROXY_PORT")
	if err != nil {
		return port, err
	} else if len(port) != 0 {
		logging.Warn("You are using the old environment variable PROXY_PORT, please consider changing it to TINYPROXY_PORT")
	} else {
		port, err = envParams.GetEnv("TINYPROXY_PORT", libparams.Default("8888"))
		if err != nil {
			return port, err
		}
	}
	return port, verification.VerifyPort(port)
}

// GetTinyProxyUser obtains the TinyProxy server user from the environment variable
// TINYPROXY_USER, and using PROXY_USER as a retro-compatibility name
func GetTinyProxyUser(envParams libparams.EnvParams) (user string, err error) {
	// Retro-compatibility
	user, err = envParams.GetEnv("PROXY_USER")
	if err != nil {
		return user, err
	}
	if len(user) != 0 {
		logging.Warn("You are using the old environment variable PROXY_USER, please consider changing it to TINYPROXY_USER")
		return user, nil
	}
	return envParams.GetEnv("TINYPROXY_USER")
}

// GetTinyProxyPassword obtains the TinyProxy server password from the environment variable
// TINYPROXY_PASSWORD, and using PROXY_PASSWORD as a retro-compatibility name
func GetTinyProxyPassword(envParams libparams.EnvParams) (password string, err error) {
	// Retro-compatibility
	password, err = envParams.GetEnv("PROXY_PASSWORD")
	if err != nil {
		return password, err
	}
	if len(password) != 0 {
		logging.Warn("You are using the old environment variable PROXY_PASSWORD, please consider changing it to TINYPROXY_PASSWORD")
		return password, nil
	}
	return envParams.GetEnv("TINYPROXY_PASSWORD")
}
