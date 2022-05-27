package env

import (
	"errors"
	"fmt"
	"strings"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/log"
)

func readLog() (log settings.Log, err error) {
	log.Level, err = readLogLevel()
	if err != nil {
		return log, err
	}

	return log, nil
}

func readLogLevel() (level *log.Level, err error) {
	s := getCleanedEnv("LOG_LEVEL")
	if s == "" {
		return nil, nil //nolint:nilnil
	}

	level = new(log.Level)
	*level, err = parseLogLevel(s)
	if err != nil {
		return nil, fmt.Errorf("environment variable LOG_LEVEL: %w", err)
	}

	return level, nil
}

var ErrLogLevelUnknown = errors.New("log level is unknown")

func parseLogLevel(s string) (level log.Level, err error) {
	switch strings.ToLower(s) {
	case "debug":
		return log.LevelDebug, nil
	case "info":
		return log.LevelInfo, nil
	case "warning":
		return log.LevelWarn, nil
	case "error":
		return log.LevelError, nil
	default:
		return level, fmt.Errorf(
			"%w: %q is not valid and can be one of debug, info, warning or error",
			ErrLogLevelUnknown, s)
	}
}
