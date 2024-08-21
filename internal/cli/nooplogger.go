package cli

type noopLogger struct{}

func newNoopLogger() *noopLogger {
	return new(noopLogger)
}

func (l *noopLogger) Info(string) {}
