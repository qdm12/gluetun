package icmp

type Logger interface {
	Debugf(format string, args ...any)
	Warnf(format string, args ...any)
}
