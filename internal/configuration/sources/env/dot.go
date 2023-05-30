package env

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gosettings/sources/env"
)

func (s *Source) readDoT() (dot settings.DoT, err error) {
	dot.Enabled, err = env.BoolPtr("DOT")
	if err != nil {
		return dot, err
	}

	dot.UpdatePeriod, err = env.DurationPtr("DNS_UPDATE_PERIOD")
	if err != nil {
		return dot, err
	}

	dot.Unbound, err = readUnbound()
	if err != nil {
		return dot, err
	}

	dot.Blacklist, err = s.readDNSBlacklist()
	if err != nil {
		return dot, err
	}

	return dot, nil
}
