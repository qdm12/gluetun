package env

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func readUnbound() (unbound settings.Unbound, err error) {
	unbound.Providers = envToCSV("DOT_PROVIDERS")

	unbound.Caching, err = envToBoolPtr("DOT_CACHING")
	if err != nil {
		return unbound, fmt.Errorf("environment variable DOT_CACHING: %w", err)
	}

	unbound.IPv6, err = envToBoolPtr("DOT_IPV6")
	if err != nil {
		return unbound, fmt.Errorf("environment variable DOT_IPV6: %w", err)
	}

	unbound.VerbosityLevel, err = envToUint8Ptr("DOT_VERBOSITY")
	if err != nil {
		return unbound, fmt.Errorf("environment variable DOT_VERBOSITY: %w", err)
	}

	unbound.VerbosityDetailsLevel, err = envToUint8Ptr("DOT_VERBOSITY_DETAILS")
	if err != nil {
		return unbound, fmt.Errorf("environment variable DOT_VERBOSITY_DETAILS: %w", err)
	}

	unbound.ValidationLogLevel, err = envToUint8Ptr("DOT_VALIDATION_LOGLEVEL")
	if err != nil {
		return unbound, fmt.Errorf("environment variable DOT_VALIDATION_LOGLEVEL: %w", err)
	}

	return unbound, nil
}
