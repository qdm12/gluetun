package configuration

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/qdm12/golibs/params"
)

// ControlServer contains settings to customize the control server operation.
type ControlServer struct {
	Port uint16
	Log  bool
}

func (settings *ControlServer) String() string {
	return strings.Join(settings.lines(), "\n")
}

func (settings *ControlServer) lines() (lines []string) {
	lines = append(lines, lastIndent+"HTTP control server:")

	lines = append(lines, indent+lastIndent+"Listening port: "+strconv.Itoa(int(settings.Port)))

	if settings.Log {
		lines = append(lines, indent+lastIndent+"Logging: enabled")
	}

	return lines
}

func (settings *ControlServer) read(r reader) (err error) {
	settings.Log, err = r.env.OnOff("HTTP_CONTROL_SERVER_LOG", params.Default("on"))
	if err != nil {
		return fmt.Errorf("environment variable HTTP_CONTROL_SERVER_LOG: %w", err)
	}

	var warning string
	settings.Port, warning, err = r.env.ListeningPort(
		"HTTP_CONTROL_SERVER_PORT", params.Default("8000"))
	if len(warning) > 0 {
		r.warner.Warn(warning)
	}
	if err != nil {
		return fmt.Errorf("environment variable HTTP_CONTROL_SERVER_PORT: %w", err)
	}

	return nil
}
