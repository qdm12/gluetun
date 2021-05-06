package resolver

import "net"

func uniqueIPsToSlice(uniqueIPs map[string]struct{}) (ips []net.IP) {
	ips = make([]net.IP, 0, len(uniqueIPs))
	for key := range uniqueIPs {
		IP := net.ParseIP(key)
		if IPv4 := IP.To4(); IPv4 != nil {
			IP = IPv4
		}
		ips = append(ips, IP)
	}
	return ips
}
