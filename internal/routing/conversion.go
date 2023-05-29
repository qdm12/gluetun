package routing

import (
	"fmt"
	"net"
	"net/netip"
)

func netIPToNetipAddress(ip net.IP) (address netip.Addr) {
	address, ok := netip.AddrFromSlice(ip)
	if !ok {
		panic(fmt.Sprintf("converting %#v to netip.Addr failed", ip))
	}
	return address.Unmap()
}
