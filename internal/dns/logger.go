package dns

type Logger interface {
	Debug(s string)
	Info(s string)
	Warn(s string)
	Error(s string)
}
