package settings

import (
	"errors"
	"fmt"
	"net/netip"
	"time"

	"github.com/qdm12/dns/v2/pkg/provider"
	"github.com/qdm12/gosettings"
	"github.com/qdm12/gosettings/reader"
	"github.com/qdm12/gotree"
)

// DoT contains settings to configure the DoT server.
type DoT struct {
	// Enabled is true if the DoT server should be running
	// and used. It defaults to true, and cannot be nil
	// in the internal state.
	Enabled *bool
	// UpdatePeriod is the period to update DNS block lists.
	// It can be set to 0 to disable the update.
	// It defaults to 24h and cannot be nil in
	// the internal state.
	UpdatePeriod *time.Duration
	// Providers is a list of DNS over TLS providers
	Providers []string `json:"providers"`
	// Caching is true if the DoT server should cache
	// DNS responses.
	Caching *bool `json:"caching"`
	// IPv6 is true if the DoT server should connect over IPv6.
	IPv6 *bool `json:"ipv6"`
	// Blacklist contains settings to configure the filter
	// block lists.
	Blacklist DNSBlacklist
}

var ErrDoTUpdatePeriodTooShort = errors.New("update period is too short")

func (d DoT) validate() (err error) {
	const minUpdatePeriod = 30 * time.Second
	if *d.UpdatePeriod != 0 && *d.UpdatePeriod < minUpdatePeriod {
		return fmt.Errorf("%w: %s must be bigger than %s",
			ErrDoTUpdatePeriodTooShort, *d.UpdatePeriod, minUpdatePeriod)
	}

	providers := provider.NewProviders()
	for _, providerName := range d.Providers {
		_, err := providers.Get(providerName)
		if err != nil {
			return err
		}
	}

	err = d.Blacklist.validate()
	if err != nil {
		return err
	}

	return nil
}

func (d *DoT) copy() (copied DoT) {
	return DoT{
		Enabled:      gosettings.CopyPointer(d.Enabled),
		UpdatePeriod: gosettings.CopyPointer(d.UpdatePeriod),
		Providers:    gosettings.CopySlice(d.Providers),
		Caching:      gosettings.CopyPointer(d.Caching),
		IPv6:         gosettings.CopyPointer(d.IPv6),
		Blacklist:    d.Blacklist.copy(),
	}
}

// overrideWith overrides fields of the receiver
// settings object with any field set in the other
// settings.
func (d *DoT) overrideWith(other DoT) {
	d.Enabled = gosettings.OverrideWithPointer(d.Enabled, other.Enabled)
	d.UpdatePeriod = gosettings.OverrideWithPointer(d.UpdatePeriod, other.UpdatePeriod)
	d.Providers = gosettings.OverrideWithSlice(d.Providers, other.Providers)
	d.Caching = gosettings.OverrideWithPointer(d.Caching, other.Caching)
	d.IPv6 = gosettings.OverrideWithPointer(d.IPv6, other.IPv6)
	d.Blacklist.overrideWith(other.Blacklist)
}

func (d *DoT) setDefaults() {
	d.Enabled = gosettings.DefaultPointer(d.Enabled, true)
	const defaultUpdatePeriod = 24 * time.Hour
	d.UpdatePeriod = gosettings.DefaultPointer(d.UpdatePeriod, defaultUpdatePeriod)
	d.Providers = gosettings.DefaultSlice(d.Providers, []string{
		provider.Cloudflare().Name,
	})
	d.Caching = gosettings.DefaultPointer(d.Caching, true)
	d.IPv6 = gosettings.DefaultPointer(d.IPv6, false)
	d.Blacklist.setDefaults()
}

func (d DoT) GetFirstPlaintextIPv4() (ipv4 netip.Addr) {
	providers := provider.NewProviders()
	provider, err := providers.Get(d.Providers[0])
	if err != nil {
		// Settings should be validated before calling this function,
		// so an error happening here is a programming error.
		panic(err)
	}

	return provider.DoT.IPv4[0].Addr()
}

func (d DoT) String() string {
	return d.toLinesNode().String()
}

func (d DoT) toLinesNode() (node *gotree.Node) {
	node = gotree.New("DNS over TLS settings:")

	node.Appendf("Enabled: %s", gosettings.BoolToYesNo(d.Enabled))
	if !*d.Enabled {
		return node
	}

	update := "disabled" //nolint:goconst
	if *d.UpdatePeriod > 0 {
		update = "every " + d.UpdatePeriod.String()
	}
	node.Appendf("Update period: %s", update)

	upstreamResolvers := node.Append("Upstream resolvers:")
	for _, provider := range d.Providers {
		upstreamResolvers.Append(provider)
	}

	node.Appendf("Caching: %s", gosettings.BoolToYesNo(d.Caching))
	node.Appendf("IPv6: %s", gosettings.BoolToYesNo(d.IPv6))

	node.AppendNode(d.Blacklist.toLinesNode())

	return node
}

func (d *DoT) read(reader *reader.Reader) (err error) {
	d.Enabled, err = reader.BoolPtr("DOT")
	if err != nil {
		return err
	}

	d.UpdatePeriod, err = reader.DurationPtr("DNS_UPDATE_PERIOD")
	if err != nil {
		return err
	}

	d.Providers = reader.CSV("DOT_PROVIDERS")

	d.Caching, err = reader.BoolPtr("DOT_CACHING")
	if err != nil {
		return err
	}

	d.IPv6, err = reader.BoolPtr("DOT_IPV6")
	if err != nil {
		return err
	}

	err = d.Blacklist.read(reader)
	if err != nil {
		return err
	}

	return nil
}
