package server

type Logger interface {
	Debugf(format string, args ...any)
	infoer
	warner
	Warnf(format string, args ...any)
	errorer
}

type infoWarner interface {
	infoer
	warner
}

type infoer interface {
	Info(s string)
}

type warner interface {
	Warn(s string)
}

type errorer interface {
	Error(s string)
}
