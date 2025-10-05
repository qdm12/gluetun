package healthcheck

type Logger interface {
	Info(s string)
	Warnf(format string, args ...any)
	Error(s string)
}
