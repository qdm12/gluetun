package firewall

type Logger interface {
	Debug(s string)
	Info(s string)
	Error(s string)
}
