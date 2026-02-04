package netlink

import (
	"fmt"
	"net"
	"net/netip"
)

func netipPrefixToIPNet(prefix netip.Prefix) (ipNet *net.IPNet) {
	if !prefix.IsValid() {
		return nil
	}

	prefixAddr := prefix.Addr().Unmap()
	ipMask := net.CIDRMask(prefix.Bits(), prefixAddr.BitLen())
	ip := netipAddrToNetIP(prefixAddr)

	return &net.IPNet{
		IP:   ip,
		Mask: ipMask,
	}
}

func netIPNetToNetipPrefix(ipNet *net.IPNet) (prefix netip.Prefix) {
	if ipNet == nil || (len(ipNet.IP) != net.IPv4len && len(ipNet.IP) != net.IPv6len) {
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

func ipAndLengthToPrefix(ip *net.IP, length uint8) netip.Prefix {
	if ip == nil || len(*ip) == 0 {
		return netip.Prefix{}
	}
	var dstIP netip.Addr
	if ipv4 := ip.To4(); ipv4 != nil { // IPv6
		dstIP = netip.AddrFrom4([4]byte(*ip))
	} else {
		dstIP = netip.AddrFrom16([16]byte(*ip))
	}
	return netip.PrefixFrom(dstIP, int(length))
}

func prefixToIPAndLength(prefix netip.Prefix) (ip *net.IP, length uint8) {
	if !prefix.IsValid() {
		return nil, 0
	}
	prefixIP := prefix.Addr().Unmap()
	ip = new(net.IP)
	*ip = netipAddrToNetIP(prefixIP)
	length = uint8(prefix.Bits()) //nolint:gosec
	return ip, length
}

func netipAddrToNetIP(address netip.Addr) (ip net.IP) {
	switch {
	case !address.IsValid():
		return nil
	case address.Is4() || address.Is4In6():
		bytes := address.As4()
		return net.IP(bytes[:])
	default:
		bytes := address.As16()
		return net.IP(bytes[:])
	}
}

func netIPToNetipAddress(ip net.IP) (address netip.Addr) {
	if len(ip) != net.IPv4len && len(ip) != net.IPv6len {
		return address // invalid
	}

	address, ok := netip.AddrFromSlice(ip)
	if !ok {
		panic(fmt.Sprintf("converting %#v to netip.Addr failed", ip))
	}
	return address.Unmap()
}
