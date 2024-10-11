//go:build !linux && !darwin

package netlink

func (n *NetLink) RouteList(family int) (
	routes []Route, err error,
) {
	panic("not implemented")
}

func (n *NetLink) RouteAdd(route Route) error {
	panic("not implemented")
}

func (n *NetLink) RouteDel(route Route) error {
	panic("not implemented")
}

func (n *NetLink) RouteReplace(route Route) error {
	panic("not implemented")
}
