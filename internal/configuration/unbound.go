package configuration

import (
	"errors"
	"fmt"
	"net"
	"strings"

	unbound "github.com/qdm12/dns/pkg/unbound"
	"github.com/qdm12/golibs/params"
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

	verbosityLevel, err := r.env.IntRange("DOT_VERBOSITY", 0, 5, params.Default("1"))
	if err != nil {
		return err
	}
	settings.Unbound.VerbosityLevel = uint8(verbosityLevel)

	verbosityDetailsLevel, err := r.env.IntRange("DOT_VERBOSITY_DETAILS", 0, 4, params.Default("0"))
	if err != nil {
		return err
	}
	settings.Unbound.VerbosityDetailsLevel = uint8(verbosityDetailsLevel)

	validationLogLevel, err := r.env.IntRange("DOT_VALIDATION_LOGLEVEL", 0, 2, params.Default("0"))
	if err != nil {
		return err
	}
	settings.Unbound.ValidationLogLevel = uint8(validationLogLevel)

	if err := settings.readUnboundPrivateAddresses(r.env); err != nil {
		return err
	}

	if err := settings.readUnboundUnblockedHostnames(r); err != nil {
		return err
	}

	settings.Unbound.AccessControl.Allowed = []net.IPNet{
		{
			IP:   net.IPv4zero,
			Mask: net.IPv4Mask(0, 0, 0, 0),
		},
		{
			IP:   net.IPv6zero,
			Mask: net.IPMask{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
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
	for _, provider := range strings.Split(s, ",") {
		_, ok := unbound.GetProviderData(provider)
		if !ok {
			return fmt.Errorf("%w: %s", ErrInvalidDNSOverTLSProvider, provider)
		}
		settings.Unbound.Providers = append(settings.Unbound.Providers, provider)
	}
	return nil
}

var (
	ErrInvalidPrivateAddress = errors.New("private address is not a valid IP or CIDR range")
)

func (settings *DNS) readUnboundPrivateAddresses(env params.Env) (err error) {
	privateAddresses, err := env.CSV("DOT_PRIVATE_ADDRESS")
	if err != nil {
		return err
	} else if len(privateAddresses) == 0 {
		return nil
	}
	for _, address := range privateAddresses {
		ip := net.ParseIP(address)
		_, _, err := net.ParseCIDR(address)
		if ip == nil && err != nil {
			return fmt.Errorf("%w: %s", ErrInvalidPrivateAddress, address)
		}
	}
	settings.Unbound.BlockedIPs = append(
		settings.Unbound.BlockedIPs, privateAddresses...)
	return nil
}

var (
	ErrInvalidHostname = errors.New("invalid hostname")
)

func (settings *DNS) readUnboundUnblockedHostnames(r reader) (err error) {
	hostnames, err := r.env.CSV("UNBLOCK")
	if err != nil {
		return err
	} else if len(hostnames) == 0 {
		return nil
	}
	for _, hostname := range hostnames {
		if !r.regex.MatchHostname(hostname) {
			return fmt.Errorf("%w: %s", ErrInvalidHostname, hostname)
		}
	}
	settings.Unbound.AllowedHostnames = append(
		settings.Unbound.AllowedHostnames, hostnames...)
	return nil
}
