package netlink

import (
	"fmt"

	"github.com/vishvananda/netlink"
)

func (n *NetLink) IsIPv6Supported() (supported bool, err error) {
	links, err := n.LinkList()
	if err != nil {
		return false, fmt.Errorf("listing links: %w", err)
	}

	for _, link := range links {
		routes, err := n.RouteList(link, netlink.FAMILY_V6)
		if err != nil {
			return false, fmt.Errorf("listing IPv6 routes for link %s: %w",
				link.Attrs().Name, err)
		}

		// Check each route for IPv6 due to Podman bug listing IPv4 routes
		// as IPv6 routes at container start, see:
		// https://github.com/qdm12/gluetun/issues/1241#issuecomment-1333405949
		for _, route := range routes {
			if route.Dst.IP.To4() == nil ||
				route.Src.To4() == nil { // Destination or source IP is IPv6
				return true, nil
			}
		}
	}

	return false, nil
}
