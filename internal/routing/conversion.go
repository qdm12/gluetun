package routing

import (
	"fmt"
	"net"
	"net/netip"
)

func NetipPrefixToIPNet(prefix *netip.Prefix) (ipNet *net.IPNet) {
	if prefix == nil {
		return nil
	}

	s := prefix.String()
	ip, ipNet, err := net.ParseCIDR(s)
	if err != nil {
		panic(err)
	}
	ipNet.IP = ip
	return ipNet
}

func netIPNetToNetipPrefix(ipNet net.IPNet) (prefix netip.Prefix) {
	return netip.MustParsePrefix(ipNet.String())
}

func netIPToNetipAddress(ip net.IP) (address netip.Addr) {
	address, ok := netip.AddrFromSlice(ip)
	if !ok {
		panic(fmt.Sprintf("converting %#v to netip.Addr failed", ip))
	}
	return address
}
