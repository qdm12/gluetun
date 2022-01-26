package httpserver

// Logger is the logger interface accepted by the
// HTTP server.
type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}
