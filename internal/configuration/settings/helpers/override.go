package helpers

import (
	"net"
	"time"

	"github.com/qdm12/golibs/logging"
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

func OverrideStringSlices(existing, other []string) (result []string) {
	if other == nil {
		return existing
	}
	result = make([]string, len(other))
	copy(result, other)
	return result
}

func OverrideUint16Slices(existing, other []uint16) (result []uint16) {
	if other == nil {
		return existing
	}
	result = make([]uint16, len(other))
	copy(result, other)
	return result
}

func OverrideIPNetsSlices(existing, other []*net.IPNet) (result []*net.IPNet) {
	if other == nil {
		return existing
	}
	result = make([]*net.IPNet, len(other))
	copy(result, other)
	return result
}
