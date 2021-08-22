package netlink

import "github.com/vishvananda/netlink"

func (n *NetLink) RouteAdd(route *netlink.Route) error {
	return netlink.RouteAdd(route)
}
