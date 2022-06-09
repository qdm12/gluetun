package updater

import (
	"fmt"
	"net"
)

func parseIPv4(s string) (ipv4 net.IP, err error) {
	ip := net.ParseIP(s)
	if ip == nil {
		return nil, fmt.Errorf("%w: %q", ErrParseIP, s)
	} else if ip.To4() == nil {
		return nil, fmt.Errorf("%w: %s", ErrNotIPv4, ip)
	}
	return ip, nil
}
