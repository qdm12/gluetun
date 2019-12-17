package settings

import "github.com/qdm12/private-internet-access-docker/internal/constants"

// TinyProxy contains settings to configure TinyProxy
type TinyProxy struct {
	Enabled  bool
	User     string
	Password string
	Port     int
	LogLevel constants.TinyProxyLogLevel
}
