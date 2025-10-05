package healthcheck

type Logger interface {
	DebugLogger
	Info(s string)
	Warnf(format string, args ...any)
	Error(s string)
}

type DebugLogger interface {
	Debug(s string)
}
