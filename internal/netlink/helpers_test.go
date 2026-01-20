package netlink

import (
	"net/netip"

	"github.com/qdm12/log"
)

func makeNetipPrefix(n byte) netip.Prefix {
	const bits = 24
	return netip.PrefixFrom(netip.AddrFrom4([4]byte{n, n, n, 0}), bits)
}

type noopLogger struct{}

func (l *noopLogger) Debug(_ string)            {}
func (l *noopLogger) Debugf(_ string, _ ...any) {}
func (l *noopLogger) Patch(_ ...log.Option)     {}
