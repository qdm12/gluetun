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

	for _, link := range links {
		routes, err := w.netlink.RouteList(link, netlink.FAMILY_V6)
		if err != nil {
			return false, fmt.Errorf("cannot list routes: %w", err)
		}

		if len(routes) > 0 {
			return true, nil
		}
	}

	return false, nil
}
