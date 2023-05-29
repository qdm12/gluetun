package wireguard

import (
	"fmt"
	"net/netip"

	"github.com/qdm12/gluetun/internal/netlink"
)

// TODO add IPv6 route if IPv6 is supported

func (w *Wireguard) addRoute(link netlink.Link, dst netip.Prefix,
	firewallMark int) (err error) {
	route := netlink.Route{
		LinkIndex: link.Index,
		Dst:       dst,
		Table:     firewallMark,
	}

	err = w.netlink.RouteAdd(route)
	if err != nil {
		return fmt.Errorf(
			"adding route for link %s, destination %s and table %d: %w",
			link.Name, dst, firewallMark, err)
	}

	return err
}
