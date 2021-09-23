package healthcheck

type Logger interface {
	Info(s string)
	Error(s string)
}
