package auth

type DebugLogger interface {
	Debugf(format string, args ...any)
	Warnf(format string, args ...any)
}
