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

// DNS contains settings to configure DNS.
type DNS struct {
	// DoTEnabled is true if the DoT server should be running
	// and used. It defaults to true, and cannot be nil
	// in the internal state.
	DoTEnabled *bool
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
	// ServerAddress is the DNS server to use inside
	// the Go program and for the system.
	// It defaults to '127.0.0.1' to be used with the
	// DoT server. It cannot be the zero value in the internal
	// state.
	ServerAddress netip.Addr
	// KeepNameserver is true if the existing DNS server
	// found in /etc/resolv.conf should be used
	// Note setting this to true will likely DNS traffic
	// outside the VPN tunnel since it would go through
	// the local DNS server of your Docker/Kubernetes
	// configuration, which is likely not going through the tunnel.
	// This will also disable the DNS over TLS server and the
	// `ServerAddress` field will be ignored.
	// It defaults to false and cannot be nil in the
	// internal state.
	KeepNameserver *bool
}

var ErrDoTUpdatePeriodTooShort = errors.New("update period is too short")

func (d DNS) validate() (err error) {
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

func (d *DNS) Copy() (copied DNS) {
	return DNS{
		DoTEnabled:     gosettings.CopyPointer(d.DoTEnabled),
		UpdatePeriod:   gosettings.CopyPointer(d.UpdatePeriod),
		Providers:      gosettings.CopySlice(d.Providers),
		Caching:        gosettings.CopyPointer(d.Caching),
		IPv6:           gosettings.CopyPointer(d.IPv6),
		Blacklist:      d.Blacklist.copy(),
		ServerAddress:  d.ServerAddress,
		KeepNameserver: gosettings.CopyPointer(d.KeepNameserver),
	}
}

// overrideWith overrides fields of the receiver
// settings object with any field set in the other
// settings.
func (d *DNS) overrideWith(other DNS) {
	d.DoTEnabled = gosettings.OverrideWithPointer(d.DoTEnabled, other.DoTEnabled)
	d.UpdatePeriod = gosettings.OverrideWithPointer(d.UpdatePeriod, other.UpdatePeriod)
	d.Providers = gosettings.OverrideWithSlice(d.Providers, other.Providers)
	d.Caching = gosettings.OverrideWithPointer(d.Caching, other.Caching)
	d.IPv6 = gosettings.OverrideWithPointer(d.IPv6, other.IPv6)
	d.Blacklist.overrideWith(other.Blacklist)
	d.ServerAddress = gosettings.OverrideWithValidator(d.ServerAddress, other.ServerAddress)
	d.KeepNameserver = gosettings.OverrideWithPointer(d.KeepNameserver, other.KeepNameserver)
}

func (d *DNS) setDefaults() {
	d.DoTEnabled = gosettings.DefaultPointer(d.DoTEnabled, true)
	const defaultUpdatePeriod = 24 * time.Hour
	d.UpdatePeriod = gosettings.DefaultPointer(d.UpdatePeriod, defaultUpdatePeriod)
	d.Providers = gosettings.DefaultSlice(d.Providers, []string{
		provider.Cloudflare().Name,
	})
	d.Caching = gosettings.DefaultPointer(d.Caching, true)
	d.IPv6 = gosettings.DefaultPointer(d.IPv6, false)
	d.Blacklist.setDefaults()
	d.ServerAddress = gosettings.DefaultValidator(d.ServerAddress,
		netip.AddrFrom4([4]byte{127, 0, 0, 1}))
	d.KeepNameserver = gosettings.DefaultPointer(d.KeepNameserver, false)
}

func (d DNS) GetFirstPlaintextIPv4() (ipv4 netip.Addr) {
	localhost := netip.AddrFrom4([4]byte{127, 0, 0, 1})
	if d.ServerAddress.Compare(localhost) != 0 && d.ServerAddress.Is4() {
		return d.ServerAddress
	}

	providers := provider.NewProviders()
	provider, err := providers.Get(d.Providers[0])
	if err != nil {
		// Settings should be validated before calling this function,
		// so an error happening here is a programming error.
		panic(err)
	}

	return provider.Plain.IPv4[0].Addr()
}

func (d DNS) String() string {
	return d.toLinesNode().String()
}

func (d DNS) toLinesNode() (node *gotree.Node) {
	node = gotree.New("DNS settings:")
	node.Appendf("Keep existing nameserver(s): %s", gosettings.BoolToYesNo(d.KeepNameserver))
	if *d.KeepNameserver {
		return node
	}
	node.Appendf("DNS server address to use: %s", d.ServerAddress)

	node.Appendf("DNS over TLS forwarder enabled: %s", gosettings.BoolToYesNo(d.DoTEnabled))
	if !*d.DoTEnabled {
		return node
	}

	update := "disabled"
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

func (d *DNS) read(r *reader.Reader) (err error) {
	d.DoTEnabled, err = r.BoolPtr("DOT")
	if err != nil {
		return err
	}

	d.UpdatePeriod, err = r.DurationPtr("DNS_UPDATE_PERIOD")
	if err != nil {
		return err
	}

	d.Providers = r.CSV("DOT_PROVIDERS")

	d.Caching, err = r.BoolPtr("DOT_CACHING")
	if err != nil {
		return err
	}

	d.IPv6, err = r.BoolPtr("DOT_IPV6")
	if err != nil {
		return err
	}

	err = d.Blacklist.read(r)
	if err != nil {
		return err
	}

	d.ServerAddress, err = r.NetipAddr("DNS_ADDRESS", reader.RetroKeys("DNS_PLAINTEXT_ADDRESS"))
	if err != nil {
		return err
	}

	d.KeepNameserver, err = r.BoolPtr("DNS_KEEP_NAMESERVER")
	if err != nil {
		return err
	}

	return nil
}
