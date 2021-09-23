package openvpn

type Logger interface {
	Debug(s string)
	Infoer
	Warn(s string)
	Error(s string)
}

type Infoer interface {
	Info(s string)
}
