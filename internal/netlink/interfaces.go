package netlink

type DebugLogger interface {
	Debugf(format string, args ...any)
}
