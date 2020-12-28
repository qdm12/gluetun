package settings

import (
	"fmt"
	"strings"
	"time"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/params"
)

type PublicIP struct {
	Period     time.Duration   `json:"period"`
	IPFilepath models.Filepath `json:"ip_filepath"`
}

func getPublicIPSettings(paramsReader params.Reader) (settings PublicIP, err error) {
	settings.Period, err = paramsReader.GetPublicIPPeriod()
	if err != nil {
		return settings, err
	}
	settings.IPFilepath, err = paramsReader.GetPublicIPFilepath()
	if err != nil {
		return settings, err
	}
	return settings, nil
}

func (s *PublicIP) String() string {
	if s.Period == 0 {
		return "Public IP getter settings: disabled"
	}
	settingsList := []string{
		"Public IP getter settings:",
		fmt.Sprintf("Period: %s", s.Period),
		fmt.Sprintf("IP file: %s", s.IPFilepath),
	}
	return strings.Join(settingsList, "\n|--")
}
