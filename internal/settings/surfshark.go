package settings

import (
	"strings"

	"github.com/qdm12/private-internet-access-docker/internal/models"
	"github.com/qdm12/private-internet-access-docker/internal/params"
)

// Surfshark contains the settings to connect to a Surfshark server
type Surfshark struct {
	User     string
	Password string
	Region   models.SurfsharkRegion
}

func (s *Surfshark) String() string {
	settingsList := []string{
		"Windscribe settings:",
		"User: [redacted]",
		"Password: [redacted]",
		"Region: " + strings.Title(string(s.Region)),
	}
	return strings.Join(settingsList, "\n |--")
}

// GetSurfsharkSettings obtains Surfshark settings from environment variables using the params package.
func GetSurfsharkSettings(paramsReader params.Reader) (settings Surfshark, err error) {
	settings.User, err = paramsReader.GetUser()
	if err != nil {
		return settings, err
	}
	settings.Password, err = paramsReader.GetPassword()
	if err != nil {
		return settings, err
	}
	settings.Region, err = paramsReader.GetSurfsharkRegion()
	if err != nil {
		return settings, err
	}
	return settings, nil
}
