package cli

import "github.com/qdm12/golibs/logging"

type noopLogger struct{}

func newNoopLogger() *noopLogger {
	return new(noopLogger)
}

func (l *noopLogger) Debug(string)             {}
func (l *noopLogger) Info(string)              {}
func (l *noopLogger) Warn(string)              {}
func (l *noopLogger) Error(string)             {}
func (l *noopLogger) PatchLevel(logging.Level) {}
func (l *noopLogger) PatchPrefix(string)       {}
