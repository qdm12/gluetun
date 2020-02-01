package constants

import (
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

const (
	// TinyProxyInfoLevel is the info log level for TinyProxy
	TinyProxyInfoLevel models.TinyProxyLogLevel = "Info"
	// TinyProxyWarnLevel is the warning log level for TinyProxy
	TinyProxyWarnLevel models.TinyProxyLogLevel = "Warning"
	// TinyProxyErrorLevel is the error log level for TinyProxy
	TinyProxyErrorLevel models.TinyProxyLogLevel = "Error"
	// TinyProxyCriticalLevel is the critical log level for TinyProxy
	TinyProxyCriticalLevel models.TinyProxyLogLevel = "Critical"
)
