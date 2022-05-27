package wireguard

import (
	"fmt"
	"net"

	"github.com/qdm12/gluetun/internal/netlink"
)

// TODO add IPv6 route if IPv6 is supported

func (w *Wireguard) addRoute(link netlink.Link, dst *net.IPNet,
	firewallMark int) (err error) {
	route := &netlink.Route{
		LinkIndex: link.Attrs().Index,
		Dst:       dst,
		Table:     firewallMark,
	}

	err = w.netlink.RouteAdd(route)
	if err != nil {
		return fmt.Errorf(
			"cannot add route for link %s, destination %s and table %d: %w",
			link.Attrs().Name, dst, firewallMark, err)
	}

	return err
}
