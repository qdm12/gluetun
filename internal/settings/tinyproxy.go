package settings

import "github.com/qdm12/private-internet-access-docker/internal/constants"

type TinyProxy struct {
	Enabled  bool
	User     string
	Password string
	Port     int
	LogLevel constants.TinyProxyLogLevel
}
