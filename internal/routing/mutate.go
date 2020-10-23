package routing

import (
	"fmt"
	"net"

	"github.com/vishvananda/netlink"
)

func (r *routing) AddRouteVia(destination net.IPNet, gateway net.IP, iface string, table int) error {
	destinationStr := destination.String()
	r.logger.Info("adding route for %s", destinationStr)
	if r.debug {
		fmt.Printf("ip route add %s via %s dev %s\n", destinationStr, gateway, iface)
	}

	link, err := netlink.LinkByName(iface)
	if err != nil {
		return fmt.Errorf("cannot add route for %s: %w", destinationStr, err)
	}
	route := netlink.Route{
		Dst:       &destination,
		Gw:        gateway,
		LinkIndex: link.Attrs().Index,
		Table:     table,
	}
	if err := netlink.RouteReplace(&route); err != nil {
		return fmt.Errorf("cannot add route for %s: %w", destinationStr, err)
	}
	return nil
}

func (r *routing) DeleteRouteVia(destination net.IPNet) (err error) {
	destinationStr := destination.String()
	r.logger.Info("deleting route for %s", destinationStr)
	if r.debug {
		fmt.Printf("ip route del %s\n", destinationStr)
	}
	route := netlink.Route{
		Dst: &destination,
	}
	if err := netlink.RouteDel(&route); err != nil {
		return fmt.Errorf("cannot delete route for %s: %w", destinationStr, err)
	}
	return nil
}
