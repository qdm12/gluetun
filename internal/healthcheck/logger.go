package healthcheck

type Logger interface {
	Debug(s string)
	Info(s string)
	Warnf(format string, args ...any)
	Error(s string)
}
