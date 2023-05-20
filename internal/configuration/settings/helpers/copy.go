package helpers

import (
	"net/netip"

	"golang.org/x/exp/slices"
)

func CopyPointer[T any](original *T) (copied *T) {
	if original == nil {
		return nil
	}
	copied = new(T)
	*copied = *original
	return copied
}

func CopySlice[T string | uint16 | netip.Addr | netip.Prefix](original []T) (copied []T) {
	return slices.Clone(original)
}
