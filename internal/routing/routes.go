package routing

import (
	"fmt"
	"net"
	"strconv"

	"github.com/qdm12/gluetun/internal/netlink"
)

func (r *Routing) addRouteVia(destination net.IPNet, gateway net.IP,
	iface string, table int) error {
	destinationStr := destination.String()
	r.logger.Info("adding route for " + destinationStr)
	r.logger.Debug("ip route replace " + destinationStr +
		" via " + gateway.String() +
		" dev " + iface +
		" table " + strconv.Itoa(table))

	link, err := r.netLinker.LinkByName(iface)
	if err != nil {
		return fmt.Errorf("cannot find link for interface %s: %w", iface, err)
	}

	route := netlink.Route{
		Dst:       &destination,
		Gw:        gateway,
		LinkIndex: link.Attrs().Index,
		Table:     table,
	}
	if err := r.netLinker.RouteReplace(&route); err != nil {
		return fmt.Errorf("cannot replace route for subnet %s at interface %s: %w",
			destinationStr, iface, err)
	}

	return nil
}

func (r *Routing) deleteRouteVia(destination net.IPNet, gateway net.IP,
	iface string, table int) (err error) {
	destinationStr := destination.String()
	r.logger.Info("deleting route for " + destinationStr)
	r.logger.Debug("ip route delete " + destinationStr +
		" via " + gateway.String() +
		" dev " + iface +
		" table " + strconv.Itoa(table))

	link, err := r.netLinker.LinkByName(iface)
	if err != nil {
		return fmt.Errorf("cannot find link for interface %s: %w", iface, err)
	}

	route := netlink.Route{
		Dst:       &destination,
		Gw:        gateway,
		LinkIndex: link.Attrs().Index,
		Table:     table,
	}
	if err := r.netLinker.RouteDel(&route); err != nil {
		return fmt.Errorf("cannot delete route: for subnet %s at interface %s: %w",
			destinationStr, iface, err)
	}

	return nil
}
