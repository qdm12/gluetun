package netlink

import (
	"fmt"
	"github.com/qdm12/gluetun/internal/pinger"
)

func (n *NetLink) IsIPv6Supported() (supported bool, err error) {
	routes, err := n.RouteList(FamilyV6)
	if err != nil {
		return false, fmt.Errorf("listing IPv6 routes: %w", err)
	}

	// Check each route for IPv6 due to Podman bug listing IPv4 routes
	// as IPv6 routes at container start, see:
	// https://github.com/qdm12/gluetun/issues/1241#issuecomment-1333405949
	for _, route := range routes {
		sourceIsIPv6 := route.Src.IsValid() && route.Src.Is6()
		destinationIsIPv6 := route.Dst.IsValid() && route.Dst.Addr().Is6()
		if sourceIsIPv6 || destinationIsIPv6 {
			link, err := n.LinkByIndex(route.LinkIndex)
			if err != nil {
				return false, fmt.Errorf("finding IPv6 supported link: %w", err)
			}
			n.debugLogger.Debugf("IPv6 is supported by link %s", link.Name)
			pingSuccess, err := pinger.Ping()
			if err != nil {
				n.debugLogger.Debugf("IPv6 support exists, but IPv6 connectivity doesn't appear to work.")
				return false, fmt.Errorf("pinging IPv6 endpoint: %w", err)
			}
			if pingSuccess == true {
				return true, nil
			}
		}
	}

	n.debugLogger.Debugf("IPv6 is not supported after searching %d routes",
		len(routes))
	return false, nil
}
