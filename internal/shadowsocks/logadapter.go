package shadowsocks

import "github.com/qdm12/golibs/logging"

type logAdapter struct {
	logger  logging.Logger
	enabled bool
}

func (l *logAdapter) Info(s string) {
	if l.enabled {
		l.logger.Info(s)
	}
}

func (l *logAdapter) Debug(s string) {
	if l.enabled {
		l.logger.Debug(s)
	}
}
func (l *logAdapter) Error(s string) {
	if l.enabled {
		l.logger.Error(s)
	}
}

func adaptLogger(logger logging.Logger, enabled bool) *logAdapter {
	return &logAdapter{
		logger:  logger,
		enabled: enabled,
	}
}
