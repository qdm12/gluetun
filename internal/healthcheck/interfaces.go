package healthcheck

type Logger interface {
	Debugf(format string, args ...any)
	Info(s string)
	Warnf(format string, args ...any)
	Error(s string)
}
