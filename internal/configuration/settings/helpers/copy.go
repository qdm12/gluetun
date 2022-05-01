package helpers

import (
	"net"
	"time"

	"github.com/qdm12/log"
	"inet.af/netaddr"
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

func CopyIP(original net.IP) (copied net.IP) {
	if original == nil {
		return nil
	}
	copied = make(net.IP, len(original))
	copy(copied, original)
	return copied
}

func CopyIPNet(original net.IPNet) (copied net.IPNet) {
	if original.IP != nil {
		copied.IP = make(net.IP, len(original.IP))
		copy(copied.IP, original.IP)
	}

	if original.Mask != nil {
		copied.Mask = make(net.IPMask, len(original.Mask))
		copy(copied.Mask, original.Mask)
	}

	return copied
}

func CopyIPNetPtr(original *net.IPNet) (copied *net.IPNet) {
	if original == nil {
		return nil
	}

	copied = new(net.IPNet)
	*copied = CopyIPNet(*original)
	return copied
}

func CopyNetaddrIP(original netaddr.IP) (copied netaddr.IP) {
	b, err := original.MarshalBinary()
	if err != nil {
		panic(err)
	}

	err = copied.UnmarshalBinary(b)
	if err != nil {
		panic(err)
	}

	return copied
}

func CopyIPPrefix(original netaddr.IPPrefix) (copied netaddr.IPPrefix) {
	b, err := original.MarshalText()
	if err != nil {
		panic(err)
	}

	err = copied.UnmarshalText(b)
	if err != nil {
		panic(err)
	}

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

func CopyIPNetSlice(original []net.IPNet) (copied []net.IPNet) {
	if original == nil {
		return nil
	}

	copied = make([]net.IPNet, len(original))
	for i := range original {
		copied[i] = CopyIPNet(original[i])
	}
	return copied
}

func CopyIPPrefixSlice(original []netaddr.IPPrefix) (copied []netaddr.IPPrefix) {
	if original == nil {
		return nil
	}

	copied = make([]netaddr.IPPrefix, len(original))
	for i := range original {
		copied[i] = CopyIPPrefix(original[i])
	}
	return copied
}

func CopyNetaddrIPsSlice(original []netaddr.IP) (copied []netaddr.IP) {
	if original == nil {
		return nil
	}

	copied = make([]netaddr.IP, len(original))
	for i := range original {
		copied[i] = CopyNetaddrIP(original[i])
	}

	return copied
}
