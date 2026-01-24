package wireguard

import (
	"math/rand/v2"
	"net/netip"

	"github.com/qdm12/gluetun/internal/netlink"
)

func ptrTo[T any](x T) *T { return &x }

var rng = rand.New(rand.NewChaCha8([32]byte{})) //nolint:gosec,gochecknoglobals

func makeLinkName() string {
	const alphabet = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, 8)
	for i := range b {
		b[i] = alphabet[rng.IntN(len(alphabet))]
	}
	return "test" + string(b)
}

func rulesAreEqual(a, b netlink.Rule) bool {
	return ipPrefixesAreEqual(a.Src, b.Src) &&
		ipPrefixesAreEqual(a.Dst, b.Dst) &&
		ptrsEqual(a.Priority, b.Priority) &&
		a.Table == b.Table &&
		a.Family == b.Family &&
		a.Flags == b.Flags &&
		a.Action == b.Action &&
		ptrsEqual(a.Mark, b.Mark)
}

func ipPrefixesAreEqual(a, b netip.Prefix) bool {
	if !a.IsValid() && !b.IsValid() {
		return true
	}
	if !a.IsValid() || !b.IsValid() {
		return false
	}
	return a.Bits() == b.Bits() &&
		a.Addr().Compare(b.Addr()) == 0
}

func ptrsEqual(a, b *uint32) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}
