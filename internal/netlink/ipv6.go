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

		if len(routes) == 0 {
			continue
		}

		return true, nil
	}

	return false, nil
}
