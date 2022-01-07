package env

import (
	"fmt"
	"os"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/govalid/binary"
	"github.com/qdm12/govalid/port"
)

func readControlServer() (controlServer settings.ControlServer, err error) {
	controlServer.Log, err = readControlServerLog()
	if err != nil {
		return controlServer, err
	}

	controlServer.Port, err = readControlServerPort()
	if err != nil {
		return controlServer, err
	}

	return controlServer, nil
}

func readControlServerLog() (enabled *bool, err error) {
	s := os.Getenv("HTTP_CONTROL_SERVER_LOG")
	if s == "" {
		return nil, nil //nolint:nilnil
	}

	log, err := binary.Validate(s)
	if err != nil {
		return nil, fmt.Errorf("environment variable HTTP_CONTROL_SERVER_LOG: %w", err)
	}

	return &log, nil
}

func readControlServerPort() (p *uint16, err error) {
	s := os.Getenv("HTTP_CONTROL_SERVER_PORT")
	if s == "" {
		return nil, nil //nolint:nilnil
	}

	p = new(uint16)
	*p, err = port.Validate(s, port.OptionPortListening(os.Geteuid()))
	if err != nil {
		return nil, fmt.Errorf("environment variable HTTP_CONTROL_SERVER_PORT: %w", err)
	}

	return p, nil
}
