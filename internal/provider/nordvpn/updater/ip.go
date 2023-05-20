package updater

import (
	"fmt"
	"net/netip"
)

func parseIPv4(s string) (ipv4 netip.Addr, err error) {
	ipv4, err = netip.ParseAddr(s)
	if err != nil {
		return ipv4, err
	}
	if !ipv4.Is4() {
		return ipv4, fmt.Errorf("%w: %s", ErrNotIPv4, ipv4)
	}
	return ipv4, nil
}
