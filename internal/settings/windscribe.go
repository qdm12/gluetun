package settings

import (
	"fmt"
	"strings"

	"github.com/qdm12/private-internet-access-docker/internal/models"
	"github.com/qdm12/private-internet-access-docker/internal/params"
)

// Windscribe contains the settings to connect to a Windscribe server
type Windscribe struct {
	User     string
	Password string
	Region   models.WindscribeRegion
	Port     uint16
}

func (w *Windscribe) String() string {
	settingsList := []string{
		"Windscribe settings:",
		"User: [redacted]",
		"Password: [redacted]",
		"Region: " + string(w.Region),
		"Custom port: " + fmt.Sprintf("%d", w.Port),
	}
	return strings.Join(settingsList, "\n |--")
}

// GetWindscribeSettings obtains Windscribe settings from environment variables using the params package.
func GetWindscribeSettings(paramsReader params.Reader, protocol models.NetworkProtocol) (settings Windscribe, err error) {
	settings.User, err = paramsReader.GetUser()
	if err != nil {
		return settings, err
	}
	settings.Password, err = paramsReader.GetPassword()
	if err != nil {
		return settings, err
	}
	settings.Region, err = paramsReader.GetWindscribeRegion()
	if err != nil {
		return settings, err
	}
	settings.Port, err = paramsReader.GetWindscribePort(protocol)
	if err != nil {
		return settings, err
	}
	return settings, nil
}
