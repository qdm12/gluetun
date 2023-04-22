package helpers

import (
	"net/netip"
)

func DefaultPointer[T any](existing *T, defaultValue T) (
	result *T) {
	if existing != nil {
		return existing
	}
	result = new(T)
	*result = defaultValue
	return result
}

func DefaultString(existing string, defaultValue string) (
	result string) {
	if existing != "" {
		return existing
	}
	return defaultValue
}

func DefaultNumber[T Number](existing T, defaultValue T) ( //nolint:ireturn
	result T) {
	if existing != 0 {
		return existing
	}
	return defaultValue
}

func DefaultIP(existing netip.Addr, defaultValue netip.Addr) (
	result netip.Addr) {
	if existing.IsValid() {
		return existing
	}
	return defaultValue
}
