package constants

import "fmt"

// TinyProxyLogLevel is the log level for TinyProxy
type TinyProxyLogLevel uint8

const (
	TinyProxyInfoLevel TinyProxyLogLevel = iota
	TinyProxyWarnLevel
	TinyProxyErrorLevel
	TinyProxyCriticalLevel
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
		return 0, fmt.Errorf("TinyProxy log level %q is not valid", s)
	}
}

func (l TinyProxyLogLevel) String() string {
	switch l {
	case TinyProxyInfoLevel:
		return "Info"
	case TinyProxyWarnLevel:
		return "Warning"
	case TinyProxyErrorLevel:
		return "Error"
	case TinyProxyCriticalLevel:
		return "Critical"
	default:
		return "INVALID"
	}
}
