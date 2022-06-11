package netlink

import "github.com/vishvananda/netlink"

type Route = netlink.Route

func (n *NetLink) RouteList(link Link, family int) (
	routes []Route, err error) {
	return netlink.RouteList(link, family)
}

func (n *NetLink) RouteAdd(route *Route) error {
	return netlink.RouteAdd(route)
}

func (n *NetLink) RouteDel(route *Route) error {
	return netlink.RouteDel(route)
}

func (n *NetLink) RouteReplace(route *Route) error {
	return netlink.RouteReplace(route)
}
