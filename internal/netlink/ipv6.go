package netlink

import (
	"fmt"
)

func (n *NetLink) IsIPv6Supported() (supported bool, err error) {
	links, err := n.LinkList()
	if err != nil {
		return false, fmt.Errorf("listing links: %w", err)
	}

	var totalRoutes uint
	for _, link := range links {
		link := link
		routes, err := n.RouteList(&link, FamilyV6)
		if err != nil {
			return false, fmt.Errorf("listing IPv6 routes for link %s: %w",
				link.Name, err)
		}

		// Check each route for IPv6 due to Podman bug listing IPv4 routes
		// as IPv6 routes at container start, see:
		// https://github.com/qdm12/gluetun/issues/1241#issuecomment-1333405949
		for _, route := range routes {
			sourceIsIPv6 := route.Src.IsValid() && route.Src.Is6()
			destinationIsIPv6 := route.Dst.IsValid() && route.Dst.Addr().Is6()
			if sourceIsIPv6 || destinationIsIPv6 {
				n.debugLogger.Debugf("IPv6 is supported by link %s", link.Name)
				return true, nil
			}
			totalRoutes++
		}
	}

	n.debugLogger.Debugf("IPv6 is not supported after searching %d links and %d routes",
		len(links), totalRoutes)
	return false, nil
}
