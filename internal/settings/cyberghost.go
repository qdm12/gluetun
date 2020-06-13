package settings

import (
	"strings"

	"github.com/qdm12/private-internet-access-docker/internal/models"
	"github.com/qdm12/private-internet-access-docker/internal/params"
)

// Windscribe contains the settings to connect to a Windscribe server
type Cyberghost struct {
	User      string
	Password  string
	ClientKey string
	Group     models.CyberghostGroup
	Region    models.CyberghostRegion
}

func (c *Cyberghost) String() string {
	settingsList := []string{
		"Cyberghost settings:",
		"User: [redacted]",
		"Password: [redacted]",
		"ClientKey: [redacted]",
		"Group: " + string(c.Group),
		"Region: " + string(c.Region),
	}
	return strings.Join(settingsList, "\n |--")
}

// GetCyberghostSettings obtains Cyberghost settings from environment variables using the params package.
func GetCyberghostSettings(paramsReader params.Reader) (settings Cyberghost, err error) {
	settings.User, err = paramsReader.GetUser()
	if err != nil {
		return settings, err
	}
	settings.Password, err = paramsReader.GetPassword()
	if err != nil {
		return settings, err
	}
	settings.ClientKey, err = paramsReader.GetCyberghostClientKey()
	if err != nil {
		return settings, err
	}
	settings.Group, err = paramsReader.GetCyberghostGroup()
	if err != nil {
		return settings, err
	}
	settings.Region, err = paramsReader.GetCyberghostRegion()
	if err != nil {
		return settings, err
	}
	return settings, nil
}
