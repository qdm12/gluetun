package httpproxy

type Logger interface {
	infoErrorer
	Debug(s string)
	Warn(s string)
}

type infoErrorer interface {
	Info(s string)
	Error(s string)
}
