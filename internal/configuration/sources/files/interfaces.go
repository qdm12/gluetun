package files

type Warner interface {
	Warnf(format string, a ...interface{})
}
