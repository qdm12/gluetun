package wireguard

import (
	"net"

	"github.com/vishvananda/netlink"
)

// TODO add IPv6 route if IPv6 is supported

func addRoute(link netlink.Link, dst *net.IPNet, firewallMark int) (err error) {
	route := &netlink.Route{
		LinkIndex: link.Attrs().Index,
		Dst:       dst,
		Table:     firewallMark,
	}

	err = netlink.RouteAdd(route)

	return err
}
