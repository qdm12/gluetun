package helpers

import (
	"net/http"
	"net/netip"
)

func OverrideWithPointer[T any](existing, other *T) (result *T) {
	if other == nil {
		return existing
	}
	result = new(T)
	*result = *other
	return result
}

func OverrideWithString(existing, other string) (result string) {
	if other == "" {
		return existing
	}
	return other
}

func OverrideWithNumber[T Number](existing, other T) (result T) { //nolint:ireturn
	if other == 0 {
		return existing
	}
	return other
}

func OverrideWithIP(existing, other netip.Addr) (result netip.Addr) {
	if !other.IsValid() {
		return existing
	}
	return other
}

func OverrideWithHTTPHandler(existing, other http.Handler) (result http.Handler) {
	if other != nil {
		return other
	}
	return existing
}

func OverrideWithSlice[T any](existing, other []T) (result []T) {
	if other == nil {
		return existing
	}
	result = make([]T, len(other))
	copy(result, other)
	return result
}
