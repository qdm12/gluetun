package env

import (
	"fmt"
	"time"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gosettings/sources/env"
)

func (s *Source) readPublicIP() (publicIP settings.PublicIP, err error) {
	publicIP.Period, err = readPublicIPPeriod()
	if err != nil {
		return publicIP, err
	}

	publicIP.IPFilepath = s.readPublicIPFilepath()

	return publicIP, nil
}

func readPublicIPPeriod() (period *time.Duration, err error) {
	s := env.Get("PUBLICIP_PERIOD")
	if s == "" {
		return nil, nil //nolint:nilnil
	}

	period = new(time.Duration)
	*period, err = time.ParseDuration(s)
	if err != nil {
		return nil, fmt.Errorf("environment variable PUBLICIP_PERIOD: %w", err)
	}

	return period, nil
}

func (s *Source) readPublicIPFilepath() (filepath *string) {
	_, value := s.getEnvWithRetro("PUBLICIP_FILE",
		[]string{"IP_STATUS_FILE"}, env.ForceLowercase(false))
	if value != "" {
		return &value
	}
	return nil
}
