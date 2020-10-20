package settings

import (
	"fmt"
	"strings"

	"github.com/qdm12/gluetun/internal/params"
)

// ControlServer contains settings to customize the control server operation.
type ControlServer struct {
	Port uint16
	Log  bool
}

func (c *ControlServer) String() string {
	settingsList := []string{
		"HTTP Control server:",
		fmt.Sprintf("Listening port: %d", c.Port),
		fmt.Sprintf("Logging: %t", c.Log),
	}
	return strings.Join(settingsList, "\n |--")
}

// GetControlServerSettings obtains the HTTP control server settings from
// environment variables using the params package.
func GetControlServerSettings(paramsReader params.Reader) (settings ControlServer, err error) {
	settings.Port, err = paramsReader.GetControlServerPort()
	if err != nil {
		return settings, err
	}
	settings.Log, err = paramsReader.GetControlServerLog()
	if err != nil {
		return settings, err
	}
	return settings, nil
}
