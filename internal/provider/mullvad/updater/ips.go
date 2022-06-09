package updater

import (
	"bytes"
	"net"
	"sort"
)

func uniqueSortedIPs(ips []net.IP) []net.IP {
	uniqueIPs := make(map[string]struct{}, len(ips))
	for _, ip := range ips {
		key := ip.String()
		uniqueIPs[key] = struct{}{}
	}

	ips = make([]net.IP, 0, len(uniqueIPs))
	for key := range uniqueIPs {
		ip := net.ParseIP(key)
		if ipv4 := ip.To4(); ipv4 != nil {
			ip = ipv4
		}
		ips = append(ips, ip)
	}

	sort.Slice(ips, func(i, j int) bool {
		return bytes.Compare(ips[i], ips[j]) < 0
	})

	return ips
}
