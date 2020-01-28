package constants

// TinyProxyLogLevel is the log level for TinyProxy
type TinyProxyLogLevel string

const (
	// TinyProxyInfoLevel is the info log level for TinyProxy
	TinyProxyInfoLevel TinyProxyLogLevel = "Info"
	// TinyProxyWarnLevel is the warning log level for TinyProxy
	TinyProxyWarnLevel = "Warning"
	// TinyProxyErrorLevel is the error log level for TinyProxy
	TinyProxyErrorLevel = "Error"
	// TinyProxyCriticalLevel is the critical log level for TinyProxy
	TinyProxyCriticalLevel = "Critical"
)
