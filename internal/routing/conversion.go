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
	if len(ipNet.IP) != net.IPv4len && len(ipNet.IP) != net.IPv6len {
		return prefix
	}
	var ip netip.Addr
	if ipv4 := ipNet.IP.To4(); ipv4 != nil {
		ip = netip.AddrFrom4([4]byte(ipv4))
	} else {
		ip = netip.AddrFrom16([16]byte(ipNet.IP))
	}
	bits, _ := ipNet.Mask.Size()
	return netip.PrefixFrom(ip, bits)
}

func netIPToNetipAddress(ip net.IP) (address netip.Addr) {
	address, ok := netip.AddrFromSlice(ip)
	if !ok {
		panic(fmt.Sprintf("converting %#v to netip.Addr failed", ip))
	}
	return address.Unmap()
}
