package env

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gosettings/sources/env"
)

func readUnbound() (unbound settings.Unbound, err error) {
	unbound.Providers = env.CSV("DOT_PROVIDERS")

	unbound.Caching, err = env.BoolPtr("DOT_CACHING")
	if err != nil {
		return unbound, err
	}

	unbound.IPv6, err = env.BoolPtr("DOT_IPV6")
	if err != nil {
		return unbound, err
	}

	unbound.VerbosityLevel, err = env.Uint8Ptr("DOT_VERBOSITY")
	if err != nil {
		return unbound, err
	}

	unbound.VerbosityDetailsLevel, err = env.Uint8Ptr("DOT_VERBOSITY_DETAILS")
	if err != nil {
		return unbound, err
	}

	unbound.ValidationLogLevel, err = env.Uint8Ptr("DOT_VALIDATION_LOGLEVEL")
	if err != nil {
		return unbound, err
	}

	return unbound, nil
}
