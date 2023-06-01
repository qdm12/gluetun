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

	_, publicIP.IPFilepath = s.getEnvWithRetro("PUBLICIP_FILE",
		[]string{"IP_STATUS_FILE"}, env.ForceLowercase(false))

	return publicIP, nil
}
