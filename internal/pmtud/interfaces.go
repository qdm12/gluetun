package pmtud

type Logger interface {
	Debug(msg string, args ...any)
}
