package helpers

import (
	"net"
	"time"

	"github.com/qdm12/golibs/logging"
	"inet.af/netaddr"
)

func OverrideWithBool(existing, other *bool) (result *bool) {
	if other == nil {
		return existing
	}
	result = new(bool)
	*result = *other
	return result
}

func OverrideWithString(existing, other string) (result string) {
	if other == "" {
		return existing
	}
	return other
}

func OverrideWithStringPtr(existing, other *string) (result *string) {
	if other == nil {
		return existing
	}
	result = new(string)
	*result = *other
	return result
}

func OverrideWithInt(existing, other *int) (result *int) {
	if other == nil {
		return existing
	}
	result = new(int)
	*result = *other
	return result
}

func OverrideWithUint8(existing, other *uint8) (result *uint8) {
	if other == nil {
		return existing
	}
	result = new(uint8)
	*result = *other
	return result
}

func OverrideWithUint16(existing, other *uint16) (result *uint16) {
	if other == nil {
		return existing
	}
	result = new(uint16)
	*result = *other
	return result
}

func OverrideWithIP(existing, other net.IP) (result net.IP) {
	if other == nil {
		return existing
	}
	result = make(net.IP, len(other))
	copy(result, other)
	return result
}

func OverrideWithDuration(existing, other *time.Duration) (result *time.Duration) {
	if other == nil {
		return existing
	}
	result = new(time.Duration)
	*result = *other
	return result
}

func OverrideWithLogLevel(existing, other *logging.Level) (result *logging.Level) {
	if other == nil {
		return existing
	}
	result = new(logging.Level)
	*result = *other
	return result
}

func OverrideWithStringSlice(existing, other []string) (result []string) {
	if other == nil {
		return existing
	}
	result = make([]string, len(other))
	copy(result, other)
	return result
}

func OverrideWithUint16Slice(existing, other []uint16) (result []uint16) {
	if other == nil {
		return existing
	}
	result = make([]uint16, len(other))
	copy(result, other)
	return result
}

func OverrideWithIPNetsSlice(existing, other []net.IPNet) (result []net.IPNet) {
	if other == nil {
		return existing
	}
	result = make([]net.IPNet, len(other))
	copy(result, other)
	return result
}

func OverrideWithNetaddrIPsSlice(existing, other []netaddr.IP) (result []netaddr.IP) {
	if other == nil {
		return existing
	}
	result = make([]netaddr.IP, len(other))
	copy(result, other)
	return result
}

func OverrideWithIPPrefixesSlice(existing, other []netaddr.IPPrefix) (result []netaddr.IPPrefix) {
	if other == nil {
		return existing
	}
	result = make([]netaddr.IPPrefix, len(other))
	copy(result, other)
	return result
}
