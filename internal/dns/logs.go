package dns

import (
	"context"
	"regexp"
	"strings"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/golibs/logging"
)

func (l *Loop) collectLines(ctx context.Context, done chan<- struct{},
	stdout, stderr chan string) {
	defer close(done)

	var line string

	for {
		select {
		case <-ctx.Done():
			// Context should only be canceled after stdout and stderr are done
			// being written to.
			close(stdout)
			close(stderr)
			return
		case line = <-stderr:
		case line = <-stdout:
		}

		line, level := processLogLine(line)
		switch level {
		case logging.LevelDebug:
			l.logger.Debug(line)
		case logging.LevelInfo:
			l.logger.Info(line)
		case logging.LevelWarn:
			l.logger.Warn(line)
		case logging.LevelError:
			l.logger.Error(line)
		}
	}
}

var unboundPrefix = regexp.MustCompile(`\[[0-9]{10}\] unbound\[[0-9]+:[0|1]\] `)

func processLogLine(s string) (filtered string, level logging.Level) {
	prefix := unboundPrefix.FindString(s)
	filtered = s[len(prefix):]
	switch {
	case strings.HasPrefix(filtered, "notice: "):
		filtered = strings.TrimPrefix(filtered, "notice: ")
		level = logging.LevelInfo
	case strings.HasPrefix(filtered, "info: "):
		filtered = strings.TrimPrefix(filtered, "info: ")
		level = logging.LevelInfo
	case strings.HasPrefix(filtered, "warn: "):
		filtered = strings.TrimPrefix(filtered, "warn: ")
		level = logging.LevelWarn
	case strings.HasPrefix(filtered, "error: "):
		filtered = strings.TrimPrefix(filtered, "error: ")
		level = logging.LevelError
	default:
		level = logging.LevelInfo
	}
	filtered = constants.ColorUnbound().Sprintf(filtered)
	return filtered, level
}
