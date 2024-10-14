package netlink

import (
	"fmt"
)

type IPv6SupportLevel uint8

const (
	IPv6Unsupported = iota
	// IPv6Supported indicates the host supports IPv6 but has no access to the
	// Internet via IPv6. It is true if one IPv6 route is found and no default
	// IPv6 route is found.
	IPv6Supported
	// IPv6Internet indicates the host has access to the Internet via IPv6,
	// which is detected when a default IPv6 route is found.
	IPv6Internet
)

func (i IPv6SupportLevel) IsSupported() bool {
	return i == IPv6Supported || i == IPv6Internet
}

func (n *NetLink) FindIPv6SupportLevel() (level IPv6SupportLevel, err error) {
	routes, err := n.RouteList(FamilyV6)
	if err != nil {
		return IPv6Unsupported, fmt.Errorf("listing IPv6 routes: %w", err)
	}

	// Check each route for IPv6 due to Podman bug listing IPv4 routes
	// as IPv6 routes at container start, see:
	// https://github.com/qdm12/gluetun/issues/1241#issuecomment-1333405949
	level = IPv6Unsupported
	for _, route := range routes {
		link, err := n.LinkByIndex(route.LinkIndex)
		if err != nil {
			return IPv6Unsupported, fmt.Errorf("finding link corresponding to route: %w", err)
		}

		sourceIsIPv4 := route.Src.IsValid() && route.Src.Is4()
		destinationIsIPv4 := route.Dst.IsValid() && route.Dst.Addr().Is4()
		destinationIsIPv6 := route.Dst.IsValid() && route.Dst.Addr().Is6()
		switch {
		case sourceIsIPv4 && destinationIsIPv4,
			destinationIsIPv6 && route.Dst.Addr().IsLoopback():
		case route.Dst.Addr().IsUnspecified(): // default ipv6 route
			n.debugLogger.Debugf("IPv6 internet access is enabled on link %s", link.Name)
			return IPv6Internet, nil
		default: // non-default ipv6 route found
			n.debugLogger.Debugf("IPv6 is supported by link %s", link.Name)
			level = IPv6Supported
		}
	}

	if level == IPv6Unsupported {
		n.debugLogger.Debugf("no IPv6 route found in %d routes", len(routes))
	}
	return level, nil
}
