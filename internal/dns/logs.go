package dns

import (
	"regexp"
	"strings"
	"sync"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/golibs/logging"
)

func (l *looper) collectLines(wg *sync.WaitGroup, stdout, stderr <-chan string) {
	defer wg.Done()
	var line string
	var ok bool
	for {
		select {
		case line, ok = <-stderr:
		case line, ok = <-stdout:
		}
		if !ok {
			return
		}
		line, level := processLogLine(line)
		switch level {
		case logging.DebugLevel:
			l.logger.Debug(line)
		case logging.InfoLevel:
			l.logger.Info(line)
		case logging.WarnLevel:
			l.logger.Warn(line)
		case logging.ErrorLevel:
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
		level = logging.InfoLevel
	case strings.HasPrefix(filtered, "info: "):
		filtered = strings.TrimPrefix(filtered, "info: ")
		level = logging.InfoLevel
	case strings.HasPrefix(filtered, "warn: "):
		filtered = strings.TrimPrefix(filtered, "warn: ")
		level = logging.WarnLevel
	case strings.HasPrefix(filtered, "error: "):
		filtered = strings.TrimPrefix(filtered, "error: ")
		level = logging.ErrorLevel
	default:
		level = logging.InfoLevel
	}
	filtered = constants.ColorUnbound().Sprintf(filtered)
	return filtered, level
}
