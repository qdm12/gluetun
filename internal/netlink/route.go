package netlink

import "github.com/vishvananda/netlink"

type Route = netlink.Route

var _ Router = (*NetLink)(nil)

type Router interface {
	RouteList(link netlink.Link, family int) (
		routes []netlink.Route, err error)
	RouteAdd(route *netlink.Route) error
	RouteDel(route *netlink.Route) error
	RouteReplace(route *netlink.Route) error
}

func (n *NetLink) RouteList(link netlink.Link, family int) (
	routes []netlink.Route, err error) {
	return netlink.RouteList(link, family)
}

func (n *NetLink) RouteAdd(route *netlink.Route) error {
	return netlink.RouteAdd(route)
}

func (n *NetLink) RouteDel(route *netlink.Route) error {
	return netlink.RouteDel(route)
}

func (n *NetLink) RouteReplace(route *netlink.Route) error {
	return netlink.RouteReplace(route)
}
