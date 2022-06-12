package env

import (
	"fmt"
	"time"

	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func readUpdater() (updater settings.Updater, err error) {
	updater.Period, err = readUpdaterPeriod()
	if err != nil {
		return updater, err
	}

	updater.DNSAddress, err = readUpdaterDNSAddress()
	if err != nil {
		return updater, err
	}

	updater.MinRatio, err = envToFloat64("UPDATER_MIN_RATIO")
	if err != nil {
		return updater, fmt.Errorf("environment variable UPDATER_MIN_RATIO: %w", err)
	}

	updater.Providers = envToCSV("UPDATER_VPN_SERVICE_PROVIDERS")

	return updater, nil
}

func readUpdaterPeriod() (period *time.Duration, err error) {
	s := getCleanedEnv("UPDATER_PERIOD")
	if s == "" {
		return nil, nil //nolint:nilnil
	}
	period = new(time.Duration)
	*period, err = time.ParseDuration(s)
	if err != nil {
		return nil, fmt.Errorf("environment variable UPDATER_PERIOD: %w", err)
	}
	return period, nil
}

func readUpdaterDNSAddress() (address string, err error) {
	// TODO this is currently using Cloudflare in
	// plaintext to not be blocked by DNS over TLS by default.
	// If a plaintext address is set in the DNS settings, this one will be used.
	// use custom future encrypted DNS written in Go without blocking
	// as it's too much trouble to start another parallel unbound instance for now.
	return "", nil
}
