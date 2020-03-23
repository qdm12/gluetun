package settings

import (
	"strings"

	"github.com/qdm12/private-internet-access-docker/internal/models"
	"github.com/qdm12/private-internet-access-docker/internal/params"
)

// Windscribe contains the settings to connect to a Windscribe server
type Windscribe struct {
	User     string
	Password string
	Region   models.WindscribeRegion
}

func (w *Windscribe) String() string {
	settingsList := []string{
		"Windscribe settings:",
		"User: [redacted]",
		"Password: [redacted]",
		"Region: " + string(w.Region),
	}
	return strings.Join(settingsList, "\n |--")
}

// GetWindscribeSettings obtains Windscribe settings from environment variables using the params package.
func GetWindscribeSettings(params params.ParamsReader) (settings Windscribe, err error) {
	settings.User, err = params.GetUser()
	if err != nil {
		return settings, err
	}
	settings.Password, err = params.GetPassword()
	if err != nil {
		return settings, err
	}
	settings.Region, err = params.GetWindscribeRegion()
	if err != nil {
		return settings, err
	}
	return settings, nil
}
