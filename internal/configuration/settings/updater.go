package settings

import (
	"fmt"
	"net"
	"time"

	"github.com/qdm12/gluetun/internal/configuration/settings/helpers"
	"github.com/qdm12/gluetun/internal/constants"
)

// Updater contains settings to configure the VPN
// server information updater.
type Updater struct {
	// Period is the period for which the updater
	// should run. It can be set to 0 to disable the
	// updater. It cannot be nil in the internal state.
	Period *time.Duration `json:"period,omitempty"`
	// DNSAddress is the DNS server address to use
	// to resolve VPN server hostnames to IP addresses.
	// It cannot be nil in the internal state.
	DNSAddress net.IP `json:"dns_address,omitempty"`
	// Providers is the list of VPN service providers
	// to update server information for.
	Providers []string `json:"providers,omitempty"`
	// CLI is to precise the updater is running in CLI
	// mode. This is set automatically and cannot be set
	// by settings sources. It cannot be nil in the
	// internal state.
	CLI *bool `json:"-"`
}

func (u Updater) validate() (err error) {
	const minPeriod = time.Minute
	if *u.Period > 0 && *u.Period < minPeriod {
		return fmt.Errorf("%w: %d must be larger than %s",
			ErrUpdaterPeriodTooSmall, *u.Period, minPeriod)
	}

	for i, provider := range u.Providers {
		valid := false
		for _, validProvider := range constants.AllProviders() {
			if provider == validProvider {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("%w: %s at index %d",
				ErrVPNProviderNameNotValid, provider, i)
		}
	}

	return nil
}

func (u *Updater) copy() (copied Updater) {
	return Updater{
		Period:     helpers.CopyDurationPtr(u.Period),
		DNSAddress: helpers.CopyIP(u.DNSAddress),
		Providers:  helpers.CopyStringSlice(u.Providers),
		CLI:        u.CLI,
	}
}

// mergeWith merges the other settings into any
// unset field of the receiver settings object.
func (u *Updater) mergeWith(other Updater) {
	u.Period = helpers.MergeWithDuration(u.Period, other.Period)
	u.DNSAddress = helpers.MergeWithIP(u.DNSAddress, other.DNSAddress)
	u.Providers = helpers.MergeStringSlices(u.Providers, other.Providers)
	u.CLI = helpers.MergeWithBool(u.CLI, other.CLI)
}

// overrideWith overrides fields of the receiver
// settings object with any field set in the other
// settings.
func (u *Updater) overrideWith(other Updater) {
	u.Period = helpers.OverrideWithDuration(u.Period, other.Period)
	u.DNSAddress = helpers.OverrideWithIP(u.DNSAddress, other.DNSAddress)
	u.Providers = helpers.OverrideStringSlices(u.Providers, other.Providers)
	u.CLI = helpers.MergeWithBool(u.CLI, other.CLI)
}

func (u *Updater) setDefaults() {
	u.Period = helpers.DefaultDuration(u.Period, 0)
	u.DNSAddress = helpers.DefaultIP(u.DNSAddress, net.IPv4(1, 1, 1, 1))
	u.CLI = helpers.DefaultBool(u.CLI, false)
}
