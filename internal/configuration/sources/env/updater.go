package env

import (
	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gosettings/sources/env"
)

func readUpdater() (updater settings.Updater, err error) {
	updater.Period, err = env.DurationPtr("UPDATER_PERIOD")
	if err != nil {
		return updater, err
	}

	updater.DNSAddress, err = readUpdaterDNSAddress()
	if err != nil {
		return updater, err
	}

	updater.MinRatio, err = env.Float64("UPDATER_MIN_RATIO")
	if err != nil {
		return updater, err
	}

	updater.Providers = env.CSV("UPDATER_VPN_SERVICE_PROVIDERS")

	return updater, nil
}

func readUpdaterDNSAddress() (address string, err error) {
	// TODO this is currently using Cloudflare in
	// plaintext to not be blocked by DNS over TLS by default.
	// If a plaintext address is set in the DNS settings, this one will be used.
	// use custom future encrypted DNS written in Go without blocking
	// as it's too much trouble to start another parallel unbound instance for now.
	return "", nil
}
