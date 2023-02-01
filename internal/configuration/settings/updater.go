package settings

import (
	"fmt"
	"strings"
	"time"

	"github.com/qdm12/gluetun/internal/configuration/settings/helpers"
	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gotree"
)

// Updater contains settings to configure the VPN
// server information updater.
type Updater struct {
	// Period is the period for which the updater
	// should run. It can be set to 0 to disable the
	// updater. It cannot be nil in the internal state.
	// TODO change to value and add Enabled field.
	Period *time.Duration
	// DNSAddress is the DNS server address to use
	// to resolve VPN server hostnames to IP addresses.
	// It cannot be the empty string in the internal state.
	DNSAddress string
	// MinRatio is the minimum ratio of servers to
	// find per provider, compared to the total current
	// number of servers. It defaults to 0.8.
	MinRatio float64
	// Providers is the list of VPN service providers
	// to update server information for.
	Providers []string
}

func (u Updater) Validate() (err error) {
	const minPeriod = time.Minute
	if *u.Period > 0 && *u.Period < minPeriod {
		return fmt.Errorf("%w: %d must be larger than %s",
			ErrUpdaterPeriodTooSmall, *u.Period, minPeriod)
	}

	if u.MinRatio <= 0 || u.MinRatio > 1 {
		return fmt.Errorf("%w: %.2f must be between 0+ and 1",
			ErrMinRatioNotValid, u.MinRatio)
	}

	validProviders := providers.All()
	for _, provider := range u.Providers {
		valid := false
		for _, validProvider := range validProviders {
			if provider == validProvider {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("%w: %q can only be one of %s",
				ErrVPNProviderNameNotValid, provider, helpers.ChoicesOrString(validProviders))
		}
	}

	return nil
}

func (u *Updater) copy() (copied Updater) {
	return Updater{
		Period:     helpers.CopyDurationPtr(u.Period),
		DNSAddress: u.DNSAddress,
		MinRatio:   u.MinRatio,
		Providers:  helpers.CopySlice(u.Providers),
	}
}

// mergeWith merges the other settings into any
// unset field of the receiver settings object.
func (u *Updater) mergeWith(other Updater) {
	u.Period = helpers.MergeWithDurationPtr(u.Period, other.Period)
	u.DNSAddress = helpers.MergeWithString(u.DNSAddress, other.DNSAddress)
	u.MinRatio = helpers.MergeWithFloat64(u.MinRatio, other.MinRatio)
	u.Providers = helpers.MergeSlices(u.Providers, other.Providers)
}

// overrideWith overrides fields of the receiver
// settings object with any field set in the other
// settings.
func (u *Updater) overrideWith(other Updater) {
	u.Period = helpers.OverrideWithDurationPtr(u.Period, other.Period)
	u.DNSAddress = helpers.OverrideWithString(u.DNSAddress, other.DNSAddress)
	u.MinRatio = helpers.OverrideWithFloat64(u.MinRatio, other.MinRatio)
	u.Providers = helpers.OverrideWithSlice(u.Providers, other.Providers)
}

func (u *Updater) SetDefaults(vpnProvider string) {
	u.Period = helpers.DefaultDurationPtr(u.Period, 0)
	u.DNSAddress = helpers.DefaultString(u.DNSAddress, "1.1.1.1:53")

	if u.MinRatio == 0 {
		const defaultMinRatio = 0.8
		u.MinRatio = defaultMinRatio
	}

	if len(u.Providers) == 0 && vpnProvider != providers.Custom {
		u.Providers = []string{vpnProvider}
	}
}

func (u Updater) String() string {
	return u.toLinesNode().String()
}

func (u Updater) toLinesNode() (node *gotree.Node) {
	if *u.Period == 0 || len(u.Providers) == 0 {
		return nil
	}

	node = gotree.New("Server data updater settings:")
	node.Appendf("Update period: %s", *u.Period)
	node.Appendf("DNS address: %s", u.DNSAddress)
	node.Appendf("Minimum ratio: %.1f", u.MinRatio)
	node.Appendf("Providers to update: %s", strings.Join(u.Providers, ", "))

	return node
}
