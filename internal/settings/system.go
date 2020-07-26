package settings

import (
	"fmt"
	"strings"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/params"
)

// System contains settings to configure system related elements
type System struct {
	UID              int
	GID              int
	Timezone         string
	IPStatusFilepath models.Filepath
}

// GetSystemSettings obtains the System settings using the params functions
func GetSystemSettings(paramsReader params.Reader) (settings System, err error) {
	settings.UID, err = paramsReader.GetUID()
	if err != nil {
		return settings, err
	}
	settings.GID, err = paramsReader.GetGID()
	if err != nil {
		return settings, err
	}
	settings.Timezone, err = paramsReader.GetTimezone()
	if err != nil {
		return settings, err
	}
	settings.IPStatusFilepath, err = paramsReader.GetIPStatusFilepath()
	if err != nil {
		return settings, err
	}
	return settings, nil
}

func (s *System) String() string {
	settingsList := []string{
		"System settings:",
		fmt.Sprintf("User ID: %d", s.UID),
		fmt.Sprintf("Group ID: %d", s.GID),
		fmt.Sprintf("Timezone: %s", s.Timezone),
		fmt.Sprintf("IP Status filepath: %s", s.IPStatusFilepath),
	}
	return strings.Join(settingsList, "\n|--")
}
