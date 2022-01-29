package env

import (
	"fmt"
	"os"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/govalid/binary"
)

func (r *Reader) readControlServer() (controlServer settings.ControlServer, err error) {
	controlServer.Log, err = readControlServerLog()
	if err != nil {
		return controlServer, err
	}

	controlServer.Address = r.readControlServerAddress()

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

func (r *Reader) readControlServerAddress() (address *string) {
	key, s := r.getEnvWithRetro("HTTP_CONTROL_SERVER_ADDRESS", "HTTP_CONTROL_SERVER_PORT")
	if s == "" {
		return nil
	}

	if key == "HTTP_CONTROL_SERVER_ADDRESS" {
		return &s
	}

	address = new(string)
	*address = ":" + s
	return address
}
