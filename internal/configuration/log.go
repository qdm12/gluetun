package configuration

import (
	"fmt"

	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/params"
)

type Log struct {
	Level logging.Level `json:"level"`
}

func (settings *Log) lines() (lines []string) {
	lines = append(lines, lastIndent+"Log:")

	lines = append(lines, indent+lastIndent+"Level: "+settings.Level.String())

	return lines
}

func (settings *Log) read(env params.Interface) (err error) {
	defaultLevel := logging.LevelInfo.String()
	settings.Level, err = env.LogLevel("LOG_LEVEL", params.Default(defaultLevel))
	if err != nil {
		return fmt.Errorf("environment variable LOG_LEVEL: %w", err)
	}

	return nil
}
