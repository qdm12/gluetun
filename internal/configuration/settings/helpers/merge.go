package helpers

import (
	"fmt"
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

func MergeWithIP(existing, other netip.Addr) (result netip.Addr) {
	if existing.IsValid() {
		return existing
	} else if !other.IsValid() {
		return existing
	}

	result, ok := netip.AddrFromSlice(other.AsSlice())
	if !ok {
		panic(fmt.Sprintf("failed copying other address: %s", other))
	}
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

func MergeSlices[K string | uint16 | netip.Addr | netip.Prefix](a, b []K) (result []K) {
	if a == nil && b == nil {
		return nil
	}

	seen := make(map[K]struct{}, len(a)+len(b))
	result = make([]K, 0, len(a)+len(b))
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
