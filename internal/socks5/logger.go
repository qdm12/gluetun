package socks5

type Logger interface {
	debuger
	infoer
	errorer
}

type debuger interface {
	Debug(s string)
}

type infoer interface {
	Info(s string)
}

type errorer interface {
	Error(s string)
}
