package netlink

import "github.com/qdm12/log"

type DebugLogger interface {
	Debug(message string)
	Debugf(format string, args ...any)
	Patch(options ...log.Option)
}
