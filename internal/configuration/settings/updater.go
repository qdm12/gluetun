package settings

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/qdm12/gluetun/internal/configuration/settings/helpers"
	"github.com/qdm12/gluetun/internal/constants"
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
	// CLI is to precise the updater is running in CLI
	// mode. This is set automatically and cannot be set
	// by settings sources. It cannot be nil in the
	// internal state.
	CLI *bool
}

func (u Updater) Validate() (err error) {
	const minPeriod = time.Minute
	if *u.Period > 0 && *u.Period < minPeriod {
		return fmt.Errorf("%w: %d must be larger than %s",
			ErrUpdaterPeriodTooSmall, *u.Period, minPeriod)
	}

	for i, provider := range u.Providers {
		valid := false
		for _, validProvider := range constants.AllProviders() {
			if validProvider == constants.Custom {
				continue
			}

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
	u.Providers = helpers.OverrideWithStringSlice(u.Providers, other.Providers)
	u.CLI = helpers.MergeWithBool(u.CLI, other.CLI)
}

func (u *Updater) SetDefaults() {
	u.Period = helpers.DefaultDuration(u.Period, 0)
	u.DNSAddress = helpers.DefaultIP(u.DNSAddress, net.IPv4(1, 1, 1, 1))
	u.CLI = helpers.DefaultBool(u.CLI, false)
}

func (u Updater) String() string {
	return u.toLinesNode().String()
}

func (u Updater) toLinesNode() (node *gotree.Node) {
	if *u.Period == 0 {
		return nil
	}

	node = gotree.New("Server data updater settings:")
	node.Appendf("Update period: %s", *u.Period)
	node.Appendf("DNS address: %s", u.DNSAddress)
	node.Appendf("Providers to update: %s", strings.Join(u.Providers, ", "))

	if *u.CLI {
		node.Appendf("CLI mode: enabled")
	}

	return node
}
