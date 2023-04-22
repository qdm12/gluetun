package helpers

import (
	"net/http"
	"net/netip"
)

func MergeWithPointer[T any](existing, other *T) (result *T) {
	if existing != nil {
		return existing
	} else if other == nil {
		return nil
	}
	result = new(T)
	*result = *other
	return result
}

func MergeWithString(existing, other string) (result string) {
	if existing != "" {
		return existing
	}
	return other
}

func MergeWithNumber[T Number](existing, other T) (result T) { //nolint:ireturn
	if existing != 0 {
		return existing
	}
	return other
}

func MergeWithIP(existing, other netip.Addr) (result netip.Addr) {
	if existing.IsValid() {
		return existing
	}
	return other
}

func MergeWithHTTPHandler(existing, other http.Handler) (result http.Handler) {
	if existing != nil {
		return existing
	}
	return other
}

func MergeSlices[T comparable](a, b []T) (result []T) {
	if a == nil && b == nil {
		return nil
	}

	seen := make(map[T]struct{}, len(a)+len(b))
	result = make([]T, 0, len(a)+len(b))
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
