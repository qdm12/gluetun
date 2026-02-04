package netlink

import (
	"math/rand/v2"
	"net/netip"

	"github.com/qdm12/log"
)

func ptrTo[T any](v T) *T { return &v }

func makeNetipPrefix(n byte) netip.Prefix {
	const bits = 24
	return netip.PrefixFrom(netip.AddrFrom4([4]byte{n, n, n, 0}), bits)
}

var rng = rand.New(rand.NewChaCha8([32]byte{})) //nolint:gosec,gochecknoglobals

func makeLinkName() string {
	const alphabet = "abcdefghijklmnopqrstuvwxyz"
	name := make([]byte, 8)
	for i := range name {
		name[i] = alphabet[rng.IntN(len(alphabet))]
	}
	return "test" + string(name)
}

type noopLogger struct{}

func (l *noopLogger) Debug(_ string)            {}
func (l *noopLogger) Debugf(_ string, _ ...any) {}
func (l *noopLogger) Patch(_ ...log.Option)     {}
