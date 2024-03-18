package env

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gosettings/sources/env"
)

func (s *Source) readPublicIP() (publicIP settings.PublicIP, err error) {
	publicIP.Period, err = s.env.DurationPtr("PUBLICIP_PERIOD")
	if err != nil {
		return publicIP, err
	}

	publicIP.IPFilepath = s.env.Get("PUBLICIP_FILE",
		env.ForceLowercase(false), env.RetroKeys("IP_STATUS_FILE"))

	publicIP.API = s.env.String("PUBLICIP_API")

	publicIP.APIToken = s.env.Get("PUBLICIP_API_TOKEN")

	return publicIP, nil
}
