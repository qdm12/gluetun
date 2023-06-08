//go:build linux || darwin

package netlink

import (
	"github.com/vishvananda/netlink"
)

func (n *NetLink) RouteList(family int) (routes []Route, err error) {
	// We set the filter to netlink.RT_FILTER_TABLE so that
	// routes from all tables are listed, as long as the filter
	// table is set to 0.
	const filterMask = netlink.RT_FILTER_TABLE
	// The filter is not left to `nil` otherwise non-main tables
	// are ignored.
	filter := &netlink.Route{}

	netlinkRoutes, err := netlink.RouteListFiltered(family, filter, filterMask)
	if err != nil {
		return nil, err
	}

	routes = make([]Route, len(netlinkRoutes))
	for i := range netlinkRoutes {
		routes[i] = netlinkRouteToRoute(netlinkRoutes[i])
	}
	return routes, nil
}

func (n *NetLink) RouteAdd(route Route) error {
	netlinkRoute := routeToNetlinkRoute(route)
	return netlink.RouteAdd(&netlinkRoute)
}

func (n *NetLink) RouteDel(route Route) error {
	netlinkRoute := routeToNetlinkRoute(route)
	return netlink.RouteDel(&netlinkRoute)
}

func (n *NetLink) RouteReplace(route Route) error {
	netlinkRoute := routeToNetlinkRoute(route)
	return netlink.RouteReplace(&netlinkRoute)
}

func netlinkRouteToRoute(netlinkRoute netlink.Route) (route Route) {
	return Route{
		LinkIndex: netlinkRoute.LinkIndex,
		Dst:       netIPNetToNetipPrefix(netlinkRoute.Dst),
		Src:       netIPToNetipAddress(netlinkRoute.Src),
		Gw:        netIPToNetipAddress(netlinkRoute.Gw),
		Priority:  netlinkRoute.Priority,
		Family:    netlinkRoute.Family,
		Table:     netlinkRoute.Table,
		Type:      netlinkRoute.Type,
	}
}

func routeToNetlinkRoute(route Route) (netlinkRoute netlink.Route) {
	return netlink.Route{
		LinkIndex: route.LinkIndex,
		Dst:       netipPrefixToIPNet(route.Dst),
		Src:       netipAddrToNetIP(route.Src),
		Gw:        netipAddrToNetIP(route.Gw),
		Priority:  route.Priority,
		Family:    route.Family,
		Table:     route.Table,
		Type:      route.Type,
	}
}
