package wireguard

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/netlink"
)

func (w *Wireguard) isIPv6Supported() (supported bool, err error) {
	links, err := w.netlink.LinkList()
	if err != nil {
		return false, fmt.Errorf("cannot list links: %w", err)
	}

	w.logger.Debug("Checking for IPv6 support...")
	for _, link := range links {
		linkName := link.Attrs().Name
		routes, err := w.netlink.RouteList(link, netlink.FAMILY_V6)
		if err != nil {
			return false, fmt.Errorf("cannot list routes for link %s: %w", linkName, err)
		}

		if len(routes) == 0 {
			w.logger.Debugf("Link %s has no IPv6 route", linkName)
			continue
		}

		w.logger.Debugf("Link %s has IPv6 routes: %#v",
			linkName, routes)
		supported = true
	}

	return supported, nil
}
