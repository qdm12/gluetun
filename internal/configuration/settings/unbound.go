package settings

import (
	"errors"
	"fmt"
	"net/netip"

	"github.com/qdm12/dns/pkg/provider"
	"github.com/qdm12/dns/pkg/unbound"
	"github.com/qdm12/gosettings"
	"github.com/qdm12/gotree"
)

// Unbound is settings for the Unbound program.
type Unbound struct {
	Providers             []string       `json:"providers"`
	Caching               *bool          `json:"caching"`
	IPv6                  *bool          `json:"ipv6"`
	VerbosityLevel        *uint8         `json:"verbosity_level"`
	VerbosityDetailsLevel *uint8         `json:"verbosity_details_level"`
	ValidationLogLevel    *uint8         `json:"validation_log_level"`
	Username              string         `json:"username"`
	Allowed               []netip.Prefix `json:"allowed"`
}

func (u *Unbound) setDefaults() {
	if len(u.Providers) == 0 {
		u.Providers = []string{
			provider.Cloudflare().String(),
		}
	}

	u.Caching = gosettings.DefaultPointer(u.Caching, true)
	u.IPv6 = gosettings.DefaultPointer(u.IPv6, false)

	const defaultVerbosityLevel = 1
	u.VerbosityLevel = gosettings.DefaultPointer(u.VerbosityLevel, defaultVerbosityLevel)

	const defaultVerbosityDetailsLevel = 0
	u.VerbosityDetailsLevel = gosettings.DefaultPointer(u.VerbosityDetailsLevel, defaultVerbosityDetailsLevel)

	const defaultValidationLogLevel = 0
	u.ValidationLogLevel = gosettings.DefaultPointer(u.ValidationLogLevel, defaultValidationLogLevel)

	if u.Allowed == nil {
		u.Allowed = []netip.Prefix{
			netip.PrefixFrom(netip.AddrFrom4([4]byte{}), 0),
			netip.PrefixFrom(netip.AddrFrom16([16]byte{}), 0),
		}
	}

	u.Username = gosettings.DefaultString(u.Username, "root")
}

var (
	ErrUnboundVerbosityLevelNotValid        = errors.New("Unbound verbosity level is not valid")
	ErrUnboundVerbosityDetailsLevelNotValid = errors.New("Unbound verbosity details level is not valid")
	ErrUnboundValidationLogLevelNotValid    = errors.New("Unbound validation log level is not valid")
)

func (u Unbound) validate() (err error) {
	for _, s := range u.Providers {
		_, err := provider.Parse(s)
		if err != nil {
			return err
		}
	}

	const maxVerbosityLevel = 5
	if *u.VerbosityLevel > maxVerbosityLevel {
		return fmt.Errorf("%w: %d must be between 0 and %d",
			ErrUnboundVerbosityLevelNotValid,
			*u.VerbosityLevel,
			maxVerbosityLevel)
	}

	const maxVerbosityDetailsLevel = 4
	if *u.VerbosityDetailsLevel > maxVerbosityDetailsLevel {
		return fmt.Errorf("%w: %d must be between 0 and %d",
			ErrUnboundVerbosityDetailsLevelNotValid,
			*u.VerbosityDetailsLevel,
			maxVerbosityDetailsLevel)
	}

	const maxValidationLogLevel = 2
	if *u.ValidationLogLevel > maxValidationLogLevel {
		return fmt.Errorf("%w: %d must be between 0 and %d",
			ErrUnboundValidationLogLevelNotValid,
			*u.ValidationLogLevel, maxValidationLogLevel)
	}

	return nil
}

func (u Unbound) copy() (copied Unbound) {
	return Unbound{
		Providers:             gosettings.CopySlice(u.Providers),
		Caching:               gosettings.CopyPointer(u.Caching),
		IPv6:                  gosettings.CopyPointer(u.IPv6),
		VerbosityLevel:        gosettings.CopyPointer(u.VerbosityLevel),
		VerbosityDetailsLevel: gosettings.CopyPointer(u.VerbosityDetailsLevel),
		ValidationLogLevel:    gosettings.CopyPointer(u.ValidationLogLevel),
		Username:              u.Username,
		Allowed:               gosettings.CopySlice(u.Allowed),
	}
}

func (u *Unbound) mergeWith(other Unbound) {
	u.Providers = gosettings.MergeWithSlice(u.Providers, other.Providers)
	u.Caching = gosettings.MergeWithPointer(u.Caching, other.Caching)
	u.IPv6 = gosettings.MergeWithPointer(u.IPv6, other.IPv6)
	u.VerbosityLevel = gosettings.MergeWithPointer(u.VerbosityLevel, other.VerbosityLevel)
	u.VerbosityDetailsLevel = gosettings.MergeWithPointer(u.VerbosityDetailsLevel, other.VerbosityDetailsLevel)
	u.ValidationLogLevel = gosettings.MergeWithPointer(u.ValidationLogLevel, other.ValidationLogLevel)
	u.Username = gosettings.MergeWithString(u.Username, other.Username)
	u.Allowed = gosettings.MergeWithSlice(u.Allowed, other.Allowed)
}

func (u *Unbound) overrideWith(other Unbound) {
	u.Providers = gosettings.OverrideWithSlice(u.Providers, other.Providers)
	u.Caching = gosettings.OverrideWithPointer(u.Caching, other.Caching)
	u.IPv6 = gosettings.OverrideWithPointer(u.IPv6, other.IPv6)
	u.VerbosityLevel = gosettings.OverrideWithPointer(u.VerbosityLevel, other.VerbosityLevel)
	u.VerbosityDetailsLevel = gosettings.OverrideWithPointer(u.VerbosityDetailsLevel, other.VerbosityDetailsLevel)
	u.ValidationLogLevel = gosettings.OverrideWithPointer(u.ValidationLogLevel, other.ValidationLogLevel)
	u.Username = gosettings.OverrideWithString(u.Username, other.Username)
	u.Allowed = gosettings.OverrideWithSlice(u.Allowed, other.Allowed)
}

func (u Unbound) ToUnboundFormat() (settings unbound.Settings, err error) {
	providers := make([]provider.Provider, len(u.Providers))
	for i := range providers {
		providers[i], err = provider.Parse(u.Providers[i])
		if err != nil {
			return settings, err
		}
	}

	const port = 53

	return unbound.Settings{
		ListeningPort:         port,
		IPv4:                  true,
		Providers:             providers,
		Caching:               *u.Caching,
		IPv6:                  *u.IPv6,
		VerbosityLevel:        *u.VerbosityLevel,
		VerbosityDetailsLevel: *u.VerbosityDetailsLevel,
		ValidationLogLevel:    *u.ValidationLogLevel,
		AccessControl: unbound.AccessControlSettings{
			Allowed: netipPrefixesToNetaddrIPPrefixes(u.Allowed),
		},
		Username: u.Username,
	}, nil
}

var (
	ErrConvertingNetip = errors.New("converting net.IP to netip.Addr failed")
)

func (u Unbound) GetFirstPlaintextIPv4() (ipv4 netip.Addr, err error) {
	s := u.Providers[0]
	provider, err := provider.Parse(s)
	if err != nil {
		return ipv4, err
	}

	ip := provider.DNS().IPv4[0]
	ipv4, ok := netip.AddrFromSlice(ip)
	if !ok {
		return ipv4, fmt.Errorf("%w: for ip %s (%#v)",
			ErrConvertingNetip, ip, ip)
	}
	return ipv4.Unmap(), nil
}

func (u Unbound) String() string {
	return u.toLinesNode().String()
}

func (u Unbound) toLinesNode() (node *gotree.Node) {
	node = gotree.New("Unbound settings:")

	authServers := node.Appendf("Authoritative servers:")
	for _, provider := range u.Providers {
		authServers.Appendf(provider)
	}

	node.Appendf("Caching: %s", gosettings.BoolToYesNo(u.Caching))
	node.Appendf("IPv6: %s", gosettings.BoolToYesNo(u.IPv6))
	node.Appendf("Verbosity level: %d", *u.VerbosityLevel)
	node.Appendf("Verbosity details level: %d", *u.VerbosityDetailsLevel)
	node.Appendf("Validation log level: %d", *u.ValidationLogLevel)
	node.Appendf("System user: %s", u.Username)

	allowedNetworks := node.Appendf("Allowed networks:")
	for _, network := range u.Allowed {
		allowedNetworks.Appendf(network.String())
	}

	return node
}
