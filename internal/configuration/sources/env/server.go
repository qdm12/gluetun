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
	const currentKey = "HTTP_CONTROL_SERVER_ADDRESS"
	key := firstKeySet(s.env, "CONTROL_SERVER_ADDRESS", currentKey)
	if key == currentKey {
		return s.env.Get(key)
	}

	s.handleDeprecatedKey(key, currentKey)
	value := s.env.Get("CONTROL_SERVER_ADDRESS")
	if value == nil {
		return nil
	}
	return ptrTo(":" + *value)
}
