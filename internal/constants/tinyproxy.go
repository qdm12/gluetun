package constants

import "fmt"

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

// ParseTinyProxyLogLevel parses a string to obtain the corresponding TinyProxyLogLevel
func ParseTinyProxyLogLevel(s string) (level TinyProxyLogLevel, err error) {
	switch s {
	case "Info":
		return TinyProxyInfoLevel, nil
	case "Warning":
		return TinyProxyWarnLevel, nil
	case "Error":
		return TinyProxyErrorLevel, nil
	case "Critical":
		return TinyProxyCriticalLevel, nil
	default:
		return "", fmt.Errorf("TinyProxy log level %q is not valid", s)
	}
}
