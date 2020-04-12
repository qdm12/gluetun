package settings

import (
	"strings"

	"github.com/qdm12/private-internet-access-docker/internal/models"
	"github.com/qdm12/private-internet-access-docker/internal/params"
)

// Mullvad contains the settings to connect to a Mullvad server
type Mullvad struct {
	User    string
	Country models.MullvadCountry
	City    models.MullvadCity
	ISP     models.MullvadProvider
	Port    uint16
}

func (m *Mullvad) String() string {
	settingsList := []string{
		"Mullvad settings:",
		"User: [redacted]",
		"Country: " + string(m.Country),
		"City: " + string(m.City),
		"ISP: " + string(m.ISP),
		"Port: " + string(m.Port),
	}
	return strings.Join(settingsList, "\n |--")
}

// GetMullvadSettings obtains Mullvad settings from environment variables using the params package.
func GetMullvadSettings(paramsReader params.Reader) (settings Mullvad, err error) {
	settings.User, err = paramsReader.GetUser()
	if err != nil {
		return settings, err
	}
	// Remove spaces in user ID to simplify user's life, thanks @JeordyR
	settings.User = strings.ReplaceAll(settings.User, " ", "")
	settings.Country, err = paramsReader.GetMullvadCountry()
	if err != nil {
		return settings, err
	}
	settings.City, err = paramsReader.GetMullvadCity()
	if err != nil {
		return settings, err
	}
	settings.ISP, err = paramsReader.GetMullvadISP()
	if err != nil {
		return settings, err
	}
	settings.Port, err = paramsReader.GetMullvadPort()
	if err != nil {
		return settings, err
	}
	return settings, nil
}
