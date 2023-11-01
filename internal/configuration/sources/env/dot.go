package env

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func (s *Source) readDoT() (dot settings.DoT, err error) {
	dot.Enabled, err = s.env.BoolPtr("DOT")
	if err != nil {
		return dot, err
	}

	dot.UpdatePeriod, err = s.env.DurationPtr("DNS_UPDATE_PERIOD")
	if err != nil {
		return dot, err
	}

	dot.Providers = s.env.CSV("DOT_PROVIDERS")

	dot.Caching, err = s.env.BoolPtr("DOT_CACHING")
	if err != nil {
		return dot, err
	}

	dot.IPv6, err = s.env.BoolPtr("DOT_IPV6")
	if err != nil {
		return dot, err
	}

	dot.Blacklist, err = s.readDNSBlacklist()
	if err != nil {
		return dot, err
	}

	return dot, nil
}
