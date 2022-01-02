package env

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func (r *Reader) readDoT() (dot settings.DoT, err error) {
	dot.Enabled, err = envToBoolPtr("DOT")
	if err != nil {
		return dot, fmt.Errorf("environment variable DOT: %w", err)
	}

	dot.UpdatePeriod, err = envToDurationPtr("DNS_UPDATE_PERIOD")
	if err != nil {
		return dot, fmt.Errorf("environment variable DNS_UPDATE_PERIOD: %w", err)
	}

	dot.Unbound, err = readUnbound()
	if err != nil {
		return dot, err
	}

	dot.Blacklist, err = r.readDNSBlacklist()
	if err != nil {
		return dot, err
	}

	return dot, nil
}
