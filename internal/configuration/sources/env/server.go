package env

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func (s *Source) readControlServer() (controlServer settings.ControlServer, err error) {
	controlServer.Log, err = s.env.BoolPtr("HTTP_CONTROL_SERVER_LOG")
	if err != nil {
		return controlServer, err
	}

	controlServer.Address = s.env.Get("HTTP_CONTROL_SERVER_ADDRESS")

	return controlServer, nil
}
