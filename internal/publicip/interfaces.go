package publicip

type Logger interface {
	Info(s string)
	Warn(s string)
	Error(s string)
}
