package settings

import (
	"github.com/qdm12/dns/pkg/provider"
	"github.com/qdm12/dns/pkg/unbound"
	"inet.af/netaddr"
)

// Unbound is settings for the Unbound program.
type Unbound struct {
	Providers             []provider.Provider
	Caching               *bool
	IPv6                  *bool
	VerbosityLevel        *uint8
	VerbosityDetailsLevel *uint8
	ValidationLogLevel    *uint8
	Allowed               []netaddr.IPPrefix
}

func (u Unbound) ToUnboundFormat() (settings unbound.Settings) {
	const port = 53
	return unbound.Settings{
		ListeningPort:         port,
		IPv4:                  true,
		Providers:             u.Providers,
		Caching:               *u.Caching,
		IPv6:                  *u.IPv6,
		VerbosityLevel:        *u.VerbosityLevel,
		VerbosityDetailsLevel: *u.VerbosityDetailsLevel,
		ValidationLogLevel:    *u.ValidationLogLevel,
		AccessControl: unbound.AccessControlSettings{
			Allowed: u.Allowed,
		},
	}
}
