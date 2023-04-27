package helpers

import (
	"net"
	"net/http"
	"net/netip"
	"time"

	"github.com/qdm12/log"
)

func MergeWithBool(existing, other *bool) (result *bool) {
	if existing != nil {
		return existing
	} else if other == nil {
		return nil
	}
	result = new(bool)
	*result = *other
	return result
}

func MergeWithString(existing, other string) (result string) {
	if existing != "" {
		return existing
	}
	return other
}

func MergeWithInt(existing, other int) (result int) {
	if existing != 0 {
		return existing
	}
	return other
}

func MergeWithFloat64(existing, other float64) (result float64) {
	if existing != 0 {
		return existing
	}
	return other
}

func MergeWithStringPtr(existing, other *string) (result *string) {
	if existing != nil {
		return existing
	} else if other == nil {
		return nil
	}
	result = new(string)
	*result = *other
	return result
}

func MergeWithIntPtr(existing, other *int) (result *int) {
	if existing != nil {
		return existing
	} else if other == nil {
		return nil
	}
	result = new(int)
	*result = *other
	return result
}

func MergeWithUint8(existing, other *uint8) (result *uint8) {
	if existing != nil {
		return existing
	} else if other == nil {
		return nil
	}
	result = new(uint8)
	*result = *other
	return result
}

func MergeWithUint16(existing, other *uint16) (result *uint16) {
	if existing != nil {
		return existing
	} else if other == nil {
		return nil
	}
	result = new(uint16)
	*result = *other
	return result
}

func MergeWithUint32(existing, other *uint32) (result *uint32) {
	if existing != nil {
		return existing
	} else if other == nil {
		return nil
	}
	result = new(uint32)
	*result = *other
	return result
}

func MergeWithIP(existing, other net.IP) (result net.IP) {
	if existing != nil {
		return existing
	} else if other == nil {
		return nil
	}
	result = make(net.IP, len(other))
	copy(result, other)
	return result
}

func MergeWithDuration(existing, other time.Duration) (result time.Duration) {
	if existing != 0 {
		return existing
	}
	return other
}

func MergeWithDurationPtr(existing, other *time.Duration) (result *time.Duration) {
	if existing != nil {
		return existing
	}
	return other
}

func MergeWithLogLevel(existing, other *log.Level) (result *log.Level) {
	if existing != nil {
		return existing
	} else if other == nil {
		return nil
	}
	result = new(log.Level)
	*result = *other
	return result
}

func MergeWithHTTPHandler(existing, other http.Handler) (result http.Handler) {
	if existing != nil {
		return existing
	}
	return other
}

func MergeStringSlices(a, b []string) (result []string) {
	if a == nil && b == nil {
		return nil
	}

	seen := make(map[string]struct{}, len(a)+len(b))
	result = make([]string, 0, len(a)+len(b))
	for _, s := range a {
		if _, ok := seen[s]; ok {
			continue // duplicate
		}
		result = append(result, s)
		seen[s] = struct{}{}
	}
	for _, s := range b {
		if _, ok := seen[s]; ok {
			continue // duplicate
		}
		result = append(result, s)
		seen[s] = struct{}{}
	}
	return result
}

func MergeUint16Slices(a, b []uint16) (result []uint16) {
	if a == nil && b == nil {
		return nil
	}

	seen := make(map[uint16]struct{}, len(a)+len(b))
	result = make([]uint16, 0, len(a)+len(b))
	for _, n := range a {
		if _, ok := seen[n]; ok {
			continue // duplicate
		}
		result = append(result, n)
		seen[n] = struct{}{}
	}
	for _, n := range b {
		if _, ok := seen[n]; ok {
			continue // duplicate
		}
		result = append(result, n)
		seen[n] = struct{}{}
	}
	return result
}

func MergeNetipAddressesSlices(a, b []netip.Addr) (result []netip.Addr) {
	if a == nil && b == nil {
		return nil
	}

	seen := make(map[string]struct{}, len(a)+len(b))
	result = make([]netip.Addr, 0, len(a)+len(b))
	for _, ip := range a {
		key := ip.String()
		if _, ok := seen[key]; ok {
			continue // duplicate
		}
		result = append(result, ip)
		seen[key] = struct{}{}
	}
	for _, ip := range b {
		key := ip.String()
		if _, ok := seen[key]; ok {
			continue // duplicate
		}
		result = append(result, ip)
		seen[key] = struct{}{}
	}
	return result
}

func MergeNetipPrefixesSlices(a, b []netip.Prefix) (result []netip.Prefix) {
	if a == nil && b == nil {
		return nil
	}

	seen := make(map[string]struct{}, len(a)+len(b))
	result = make([]netip.Prefix, 0, len(a)+len(b))
	for _, ipPrefix := range a {
		key := ipPrefix.String()
		if _, ok := seen[key]; ok {
			continue // duplicate
		}
		result = append(result, ipPrefix)
		seen[key] = struct{}{}
	}
	for _, ipPrefix := range b {
		key := ipPrefix.String()
		if _, ok := seen[key]; ok {
			continue // duplicate
		}
		result = append(result, ipPrefix)
		seen[key] = struct{}{}
	}
	return result
}
