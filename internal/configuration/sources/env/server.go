package env

import (
	"fmt"
	"os"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/govalid/binary"
	"github.com/qdm12/govalid/port"
)

func (r *Reader) readControlServer() (controlServer settings.ControlServer, err error) {
	controlServer.Log, err = readControlServerLog()
	if err != nil {
		return controlServer, err
	}

	controlServer.Address, err = r.readControlServerAddress()
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

func (r *Reader) readControlServerAddress() (address *string, err error) {
	// Retro-compatibility
	s := os.Getenv("HTTP_CONTROL_SERVER_PORT")
	if s != "" {
		r.onRetroActive("HTTP_CONTROL_SERVER_PORT", "HTTP_CONTROL_SERVER_ADDRESS")
		port, err := port.Validate(s)
		if err != nil {
			return nil, fmt.Errorf("environment variable HTTP_CONTROL_SERVER_PORT: %w", err)
		}
		address = new(string)
		*address = ":" + fmt.Sprint(port)
		return address, nil
	}

	s = os.Getenv("HTTP_CONTROL_SERVER_ADDRESS")
	if s == "" {
		return nil, nil //nolint:nilnil
	}
	return &s, nil
}
