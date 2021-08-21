package wireguard

import (
	"fmt"
	"net"

	"github.com/vishvananda/netlink"
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
		return fmt.Errorf("%w: when adding route: %s", err, route)
	}

	return err
}
