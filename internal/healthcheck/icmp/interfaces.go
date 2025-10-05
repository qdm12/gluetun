package icmp

type Warner interface {
	Warnf(format string, args ...any)
}
