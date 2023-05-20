package updater

import (
	"net/netip"
	"sort"
)

func uniqueSortedIPs(ips []netip.Addr) []netip.Addr {
	uniqueIPs := make(map[string]struct{}, len(ips))
	for _, ip := range ips {
		key := ip.String()
		uniqueIPs[key] = struct{}{}
	}

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

	sort.Slice(ips, func(i, j int) bool {
		return ips[i].Compare(ips[j]) < 0
	})

	return ips
}
