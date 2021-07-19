package configuration

import (
	"errors"
	"fmt"
	"strings"

	"github.com/qdm12/dns/pkg/provider"
	"github.com/qdm12/golibs/params"
	"inet.af/netaddr"
)

func (settings *DNS) readUnbound(r reader) (err error) {
	if err := settings.readUnboundProviders(r.env); err != nil {
		return err
	}

	settings.Unbound.ListeningPort = 53

	settings.Unbound.Caching, err = r.env.OnOff("DOT_CACHING", params.Default("on"))
	if err != nil {
		return err
	}

	settings.Unbound.IPv4 = true

	settings.Unbound.IPv6, err = r.env.OnOff("DOT_IPV6", params.Default("off"))
	if err != nil {
		return err
	}

	verbosityLevel, err := r.env.IntRange("DOT_VERBOSITY", 0, 5, params.Default("1")) //nolint:gomnd
	if err != nil {
		return err
	}
	settings.Unbound.VerbosityLevel = uint8(verbosityLevel)

	verbosityDetailsLevel, err := r.env.IntRange("DOT_VERBOSITY_DETAILS", 0, 4, params.Default("0")) //nolint:gomnd
	if err != nil {
		return err
	}
	settings.Unbound.VerbosityDetailsLevel = uint8(verbosityDetailsLevel)

	validationLogLevel, err := r.env.IntRange("DOT_VALIDATION_LOGLEVEL", 0, 2, params.Default("0")) //nolint:gomnd
	if err != nil {
		return err
	}
	settings.Unbound.ValidationLogLevel = uint8(validationLogLevel)

	settings.Unbound.AccessControl.Allowed = []netaddr.IPPrefix{
		netaddr.IPPrefixFrom(netaddr.IPv4(0, 0, 0, 0), 0),
		netaddr.IPPrefixFrom(netaddr.IPv6Raw([16]byte{}), 0),
	}

	return nil
}

var (
	ErrInvalidDNSOverTLSProvider = errors.New("invalid DNS over TLS provider")
)

func (settings *DNS) readUnboundProviders(env params.Env) (err error) {
	s, err := env.Get("DOT_PROVIDERS", params.Default("cloudflare"))
	if err != nil {
		return err
	}
	for _, field := range strings.Split(s, ",") {
		dnsProvider, err := provider.Parse(field)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrInvalidDNSOverTLSProvider, err)
		}
		settings.Unbound.Providers = append(settings.Unbound.Providers, dnsProvider)
	}
	return nil
}

var (
	ErrInvalidHostname = errors.New("invalid hostname")
)
