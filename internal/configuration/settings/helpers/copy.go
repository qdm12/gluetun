package helpers

import (
	"net/netip"
	"time"

	"github.com/qdm12/log"
)

func CopyStringPtr(original *string) (copied *string) {
	if original == nil {
		return nil
	}
	copied = new(string)
	*copied = *original
	return copied
}

func CopyBoolPtr(original *bool) (copied *bool) {
	if original == nil {
		return nil
	}
	copied = new(bool)
	*copied = *original
	return copied
}

func CopyUint8Ptr(original *uint8) (copied *uint8) {
	if original == nil {
		return nil
	}
	copied = new(uint8)
	*copied = *original
	return copied
}

func CopyUint16Ptr(original *uint16) (copied *uint16) {
	if original == nil {
		return nil
	}
	copied = new(uint16)
	*copied = *original
	return copied
}

func CopyUint32Ptr(original *uint32) (copied *uint32) {
	if original == nil {
		return nil
	}
	copied = new(uint32)
	*copied = *original
	return copied
}

func CopyIntPtr(original *int) (copied *int) {
	if original == nil {
		return nil
	}
	copied = new(int)
	*copied = *original
	return copied
}

func CopyDurationPtr(original *time.Duration) (copied *time.Duration) {
	if original == nil {
		return nil
	}
	copied = new(time.Duration)
	*copied = *original
	return copied
}

func CopyLogLevelPtr(original *log.Level) (copied *log.Level) {
	if original == nil {
		return nil
	}
	copied = new(log.Level)
	*copied = *original
	return copied
}

func CopyStringSlice(original []string) (copied []string) {
	if original == nil {
		return nil
	}

	copied = make([]string, len(original))
	copy(copied, original)
	return copied
}

func CopyUint16Slice(original []uint16) (copied []uint16) {
	if original == nil {
		return nil
	}

	copied = make([]uint16, len(original))
	copy(copied, original)
	return copied
}

func CopyNetipPrefixesSlice(original []netip.Prefix) (copied []netip.Prefix) {
	if original == nil {
		return nil
	}

	copied = make([]netip.Prefix, len(original))
	copy(copied, original)
	return copied
}

func CopyNetipAddressesSlice(original []netip.Addr) (copied []netip.Addr) {
	if original == nil {
		return nil
	}

	copied = make([]netip.Addr, len(original))
	copy(copied, original)
	return copied
}
