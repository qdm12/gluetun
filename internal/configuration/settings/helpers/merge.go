package helpers

import (
	"net"
	"time"

	"github.com/qdm12/golibs/logging"
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

func MergeWithInt(existing, other *int) (result *int) {
	if existing != nil {
		return existing
	} else if other == nil {
		return nil
	}
	result = new(int)
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

func MergeWithDuration(existing, other *time.Duration) (result *time.Duration) {
	if existing != nil {
		return existing
	}
	return other
}

func MergeWithLogLevel(existing, other *logging.Level) (result *logging.Level) {
	if existing != nil {
		return existing
	} else if other == nil {
		return nil
	}
	result = new(logging.Level)
	*result = *other
	return result
}

func MergeStringSlices(a, b []string) (result []string) {
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

func MergeIPNetsSlices(a, b []*net.IPNet) (result []*net.IPNet) {
	seen := make(map[string]struct{}, len(a)+len(b))
	result = make([]*net.IPNet, 0, len(a)+len(b))
	for _, ipNet := range a {
		if ipNet == nil {
			continue
		}
		key := ipNet.String()
		if _, ok := seen[key]; ok {
			continue // duplicate
		}
		result = append(result, ipNet)
		seen[key] = struct{}{}
	}
	for _, ipNet := range b {
		if ipNet == nil {
			continue
		}
		key := ipNet.String()
		if _, ok := seen[key]; ok {
			continue // duplicate
		}
		result = append(result, ipNet)
		seen[key] = struct{}{}
	}
	return result
}
