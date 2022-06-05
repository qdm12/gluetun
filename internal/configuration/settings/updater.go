package settings

import (
	"fmt"
	"net"
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
	// It cannot be nil in the internal state.
	DNSAddress net.IP
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
		DNSAddress: helpers.CopyIP(u.DNSAddress),
		Providers:  helpers.CopyStringSlice(u.Providers),
	}
}

// mergeWith merges the other settings into any
// unset field of the receiver settings object.
func (u *Updater) mergeWith(other Updater) {
	u.Period = helpers.MergeWithDuration(u.Period, other.Period)
	u.DNSAddress = helpers.MergeWithIP(u.DNSAddress, other.DNSAddress)
	u.Providers = helpers.MergeStringSlices(u.Providers, other.Providers)
}

// overrideWith overrides fields of the receiver
// settings object with any field set in the other
// settings.
func (u *Updater) overrideWith(other Updater) {
	u.Period = helpers.OverrideWithDuration(u.Period, other.Period)
	u.DNSAddress = helpers.OverrideWithIP(u.DNSAddress, other.DNSAddress)
	u.Providers = helpers.OverrideWithStringSlice(u.Providers, other.Providers)
}

func (u *Updater) SetDefaults(vpnProvider string) {
	u.Period = helpers.DefaultDuration(u.Period, 0)
	u.DNSAddress = helpers.DefaultIP(u.DNSAddress, net.IPv4(1, 1, 1, 1))
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
	node.Appendf("Providers to update: %s", strings.Join(u.Providers, ", "))

	return node
}
