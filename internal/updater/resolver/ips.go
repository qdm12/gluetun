package resolver

import (
	"net/netip"
)

func uniqueIPsToSlice(uniqueIPs map[string]struct{}) (ips []netip.Addr) {
	ips = make([]netip.Addr, 0, len(uniqueIPs))
	for key := range uniqueIPs {
		ip, err := netip.ParseAddr(key)
		if err != nil {
			panic(err)
		}
		if ip.Is4In6() {
			ip = netip.AddrFrom4(ip.As4())
		}
		ips = append(ips, ip)
	}
	return ips
}
