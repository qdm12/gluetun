package cli

import "github.com/qdm12/golibs/logging"

type noopLogger struct{}

func newNoopLogger() *noopLogger {
	return new(noopLogger)
}

func (l *noopLogger) Debug(s string)                 {}
func (l *noopLogger) Info(s string)                  {}
func (l *noopLogger) Warn(s string)                  {}
func (l *noopLogger) Error(s string)                 {}
func (l *noopLogger) PatchLevel(level logging.Level) {}
func (l *noopLogger) PatchPrefix(prefix string)      {}
