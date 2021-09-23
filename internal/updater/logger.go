package updater

type Logger interface {
	infoer
	warner
	errorer
}

type infoErrorer interface {
	infoer
	errorer
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
