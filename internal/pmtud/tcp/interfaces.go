package tcp

type Logger interface {
	Debug(msg string)
	Debugf(msg string, args ...any)
	Warnf(msg string, args ...any)
}
