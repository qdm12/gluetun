package settings

import (
	"fmt"
	"strings"

	"github.com/qdm12/gluetun/internal/params"
)

// System contains settings to configure system related elements.
type System struct {
	PUID     int
	PGID     int
	Timezone string
}

// GetSystemSettings obtains the System settings using the params functions.
func GetSystemSettings(paramsReader params.Reader) (settings System, err error) {
	settings.PUID, err = paramsReader.GetPUID()
	if err != nil {
		return settings, err
	}
	settings.PGID, err = paramsReader.GetPGID()
	if err != nil {
		return settings, err
	}
	settings.Timezone, err = paramsReader.GetTimezone()
	if err != nil {
		return settings, err
	}
	return settings, nil
}

func (s *System) String() string {
	settingsList := []string{
		"System settings:",
		fmt.Sprintf("Process user ID: %d", s.PUID),
		fmt.Sprintf("Process group ID: %d", s.PGID),
		fmt.Sprintf("Timezone: %s", s.Timezone),
	}
	return strings.Join(settingsList, "\n|--")
}
