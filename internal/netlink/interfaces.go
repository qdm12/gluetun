package netlink

import "github.com/qdm12/log"

type DebugLogger interface {
	Debugf(format string, args ...any)
	Patch(options ...log.Option)
}
