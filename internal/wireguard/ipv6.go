package wireguard

import (
	"errors"
	"fmt"

	"github.com/qdm12/gluetun/internal/netlink"
)

var (
	errLinkList  = errors.New("cannot list links")
	errRouteList = errors.New("cannot list routes")
)

func (w *Wireguard) isIPv6Supported() (supported bool, err error) {
	links, err := w.netlink.LinkList()
	if err != nil {
		return false, fmt.Errorf("%w: %s", errLinkList, err)
	}

	for _, link := range links {
		routes, err := w.netlink.RouteList(link, netlink.FAMILY_V6)
		if err != nil {
			return false, fmt.Errorf("%w: %s", errRouteList, err)
		}

		if len(routes) > 0 {
			return true, nil
		}
	}

	return false, nil
}
