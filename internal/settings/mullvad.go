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
func GetMullvadSettings(params params.ParamsReader) (settings Mullvad, err error) {
	settings.User, err = params.GetUser()
	if err != nil {
		return settings, err
	}
	// Remove spaces in user ID to simplify user's life, thanks @JeordyR
	settings.User = strings.ReplaceAll(settings.User, " ", "")
	settings.Country, err = params.GetMullvadCountry()
	if err != nil {
		return settings, err
	}
	settings.City, err = params.GetMullvadCity()
	if err != nil {
		return settings, err
	}
	settings.ISP, err = params.GetMullvadISP()
	if err != nil {
		return settings, err
	}
	settings.Port, err = params.GetMullvadPort()
	if err != nil {
		return settings, err
	}
	return settings, nil
}
