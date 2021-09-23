package httpproxy

type Logger interface {
	Debug(s string)
	infoer
	Warn(s string)
	errorer
}

type infoErrorer interface {
	infoer
	errorer
}

type infoer interface {
	Info(s string)
}

type errorer interface {
	Error(s string)
}
