package env

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/golibs/logging"
)

func readLog() (log settings.Log, err error) {
	log.Level, err = readLogLevel()
	if err != nil {
		return log, err
	}

	return log, nil
}

func readLogLevel() (level *logging.Level, err error) {
	s := os.Getenv("LOG_LEVEL")
	if s == "" {
		return nil, nil //nolint:nilnil
	}

	level = new(logging.Level)
	*level, err = parseLogLevel(s)
	if err != nil {
		return nil, fmt.Errorf("environment variable LOG_LEVEL: %w", err)
	}

	return level, nil
}

var ErrLogLevelUnknown = errors.New("log level is unknown")

func parseLogLevel(s string) (level logging.Level, err error) {
	switch strings.ToLower(s) {
	case "debug":
		return logging.LevelDebug, nil
	case "info":
		return logging.LevelInfo, nil
	case "warning":
		return logging.LevelWarn, nil
	case "error":
		return logging.LevelError, nil
	default:
		return level, fmt.Errorf(
			"%w: %q is not valid and can be one of debug, info, warning or error",
			ErrLogLevelUnknown, s)
	}
}
