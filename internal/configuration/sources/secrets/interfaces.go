package secrets

type Warner interface {
	Warnf(format string, a ...interface{})
}
