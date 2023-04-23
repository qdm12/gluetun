package settings

import (
	"net/netip"

	"inet.af/netaddr"
)

func netipAddressToNetaddrIP(address netip.Addr) (ip netaddr.IP) {
	if address.Is4() {
		return netaddr.IPFrom4(address.As4())
	}
	return netaddr.IPFrom16(address.As16())
}

func netipAddressesToNetaddrIPs(addresses []netip.Addr) (ips []netaddr.IP) {
	ips = make([]netaddr.IP, len(addresses))
	for i := range addresses {
		ips[i] = netipAddressToNetaddrIP(addresses[i])
	}
	return ips
}

func netipPrefixToNetaddrIPPrefix(prefix netip.Prefix) (ipPrefix netaddr.IPPrefix) {
	netaddrIP := netipAddressToNetaddrIP(prefix.Addr())
	bits := prefix.Bits()
	return netaddr.IPPrefixFrom(netaddrIP, uint8(bits))
}

func netipPrefixesToNetaddrIPPrefixes(prefixes []netip.Prefix) (ipPrefixes []netaddr.IPPrefix) {
	ipPrefixes = make([]netaddr.IPPrefix, len(prefixes))
	for i := range ipPrefixes {
		ipPrefixes[i] = netipPrefixToNetaddrIPPrefix(prefixes[i])
	}
	return ipPrefixes
}
