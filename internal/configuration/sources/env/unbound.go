package env

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func (s *Source) readUnbound() (unbound settings.Unbound, err error) {
	unbound.Providers = s.env.CSV("DOT_PROVIDERS")

	unbound.Caching, err = s.env.BoolPtr("DOT_CACHING")
	if err != nil {
		return unbound, err
	}

	unbound.IPv6, err = s.env.BoolPtr("DOT_IPV6")
	if err != nil {
		return unbound, err
	}

	unbound.VerbosityLevel, err = s.env.Uint8Ptr("DOT_VERBOSITY")
	if err != nil {
		return unbound, err
	}

	unbound.VerbosityDetailsLevel, err = s.env.Uint8Ptr("DOT_VERBOSITY_DETAILS")
	if err != nil {
		return unbound, err
	}

	unbound.ValidationLogLevel, err = s.env.Uint8Ptr("DOT_VALIDATION_LOGLEVEL")
	if err != nil {
		return unbound, err
	}

	return unbound, nil
}
