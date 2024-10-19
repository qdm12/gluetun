package netlink

import "net/netip"

func makeNetipPrefix(n byte) netip.Prefix {
	const bits = 24
	return netip.PrefixFrom(netip.AddrFrom4([4]byte{n, n, n, 0}), bits)
}
