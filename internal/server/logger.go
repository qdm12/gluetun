package server

type Logger interface {
	infoer
	warner
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
