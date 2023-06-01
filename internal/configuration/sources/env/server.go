package env

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func (s *Source) readControlServer() (controlServer settings.ControlServer, err error) {
	controlServer.Log, err = s.env.BoolPtr("HTTP_CONTROL_SERVER_LOG")
	if err != nil {
		return controlServer, err
	}

	controlServer.Address = s.readControlServerAddress()

	return controlServer, nil
}

func (s *Source) readControlServerAddress() (address *string) {
	key, value := s.getEnvWithRetro("HTTP_CONTROL_SERVER_ADDRESS",
		[]string{"CONTROL_SERVER_ADDRESS"})
	if value == nil {
		return nil
	}

	if key == "HTTP_CONTROL_SERVER_ADDRESS" {
		return value
	}

	address = new(string)
	*address = ":" + *value
	return address
}
