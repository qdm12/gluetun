package settings

import (
	"errors"
	"fmt"
	"net"

	"github.com/qdm12/dns/pkg/provider"
	"github.com/qdm12/dns/pkg/unbound"
	"github.com/qdm12/gluetun/internal/configuration/settings/helpers"
	"github.com/qdm12/gotree"
	"inet.af/netaddr"
)

// Unbound is settings for the Unbound program.
type Unbound struct {
	Providers             []string
	Caching               *bool
	IPv6                  *bool
	VerbosityLevel        *uint8
	VerbosityDetailsLevel *uint8
	ValidationLogLevel    *uint8
	Username              string
	Allowed               []netaddr.IPPrefix
}

func (u *Unbound) setDefaults() {
	if len(u.Providers) == 0 {
		u.Providers = []string{
			provider.Cloudflare().String(),
		}
	}

	u.Caching = helpers.DefaultBool(u.Caching, true)
	u.IPv6 = helpers.DefaultBool(u.IPv6, false)

	const defaultVerbosityLevel = 1
	u.VerbosityLevel = helpers.DefaultUint8(u.VerbosityLevel, defaultVerbosityLevel)

	const defaultVerbosityDetailsLevel = 0
	u.VerbosityDetailsLevel = helpers.DefaultUint8(u.VerbosityDetailsLevel, defaultVerbosityDetailsLevel)

	const defaultValidationLogLevel = 0
	u.ValidationLogLevel = helpers.DefaultUint8(u.ValidationLogLevel, defaultValidationLogLevel)

	if u.Allowed == nil {
		u.Allowed = []netaddr.IPPrefix{
			netaddr.IPPrefixFrom(netaddr.IPv4(0, 0, 0, 0), 0),
			netaddr.IPPrefixFrom(netaddr.IPv6Raw([16]byte{}), 0),
		}
	}

	u.Username = helpers.DefaultString(u.Username, "root")
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
		Providers:             helpers.CopyStringSlice(u.Providers),
		Caching:               helpers.CopyBoolPtr(u.Caching),
		IPv6:                  helpers.CopyBoolPtr(u.IPv6),
		VerbosityLevel:        helpers.CopyUint8Ptr(u.VerbosityLevel),
		VerbosityDetailsLevel: helpers.CopyUint8Ptr(u.VerbosityDetailsLevel),
		ValidationLogLevel:    helpers.CopyUint8Ptr(u.ValidationLogLevel),
		Username:              u.Username,
		Allowed:               helpers.CopyIPPrefixSlice(u.Allowed),
	}
}

func (u *Unbound) mergeWith(other Unbound) {
	u.Providers = helpers.MergeStringSlices(u.Providers, other.Providers)
	u.Caching = helpers.MergeWithBool(u.Caching, other.Caching)
	u.IPv6 = helpers.MergeWithBool(u.IPv6, other.IPv6)
	u.VerbosityLevel = helpers.MergeWithUint8(u.VerbosityLevel, other.VerbosityLevel)
	u.VerbosityDetailsLevel = helpers.MergeWithUint8(u.VerbosityDetailsLevel, other.VerbosityDetailsLevel)
	u.ValidationLogLevel = helpers.MergeWithUint8(u.ValidationLogLevel, other.ValidationLogLevel)
	u.Username = helpers.MergeWithString(u.Username, other.Username)
	u.Allowed = helpers.MergeIPPrefixesSlices(u.Allowed, other.Allowed)
}

func (u *Unbound) overrideWith(other Unbound) {
	u.Providers = helpers.OverrideWithStringSlice(u.Providers, other.Providers)
	u.Caching = helpers.OverrideWithBool(u.Caching, other.Caching)
	u.IPv6 = helpers.OverrideWithBool(u.IPv6, other.IPv6)
	u.VerbosityLevel = helpers.OverrideWithUint8(u.VerbosityLevel, other.VerbosityLevel)
	u.VerbosityDetailsLevel = helpers.OverrideWithUint8(u.VerbosityDetailsLevel, other.VerbosityDetailsLevel)
	u.ValidationLogLevel = helpers.OverrideWithUint8(u.ValidationLogLevel, other.ValidationLogLevel)
	u.Username = helpers.OverrideWithString(u.Username, other.Username)
	u.Allowed = helpers.OverrideWithIPPrefixesSlice(u.Allowed, other.Allowed)
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
			Allowed: u.Allowed,
		},
		Username: u.Username,
	}, nil
}

func (u Unbound) GetFirstPlaintextIPv4() (ipv4 net.IP, err error) {
	s := u.Providers[0]
	provider, err := provider.Parse(s)
	if err != nil {
		return nil, err
	}

	return provider.DNS().IPv4[0], nil
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

	node.Appendf("Caching: %s", helpers.BoolPtrToYesNo(u.Caching))
	node.Appendf("IPv6: %s", helpers.BoolPtrToYesNo(u.IPv6))
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
