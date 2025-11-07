package healthcheck

type Logger interface {
	Debugf(format string, args ...any)
	Info(s string)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Error(s string)
}
