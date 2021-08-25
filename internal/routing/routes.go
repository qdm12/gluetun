package routing

import (
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/qdm12/gluetun/internal/netlink"
)

var (
	errLinkByName = errors.New("cannot obtain link by name")
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
		return fmt.Errorf("%w: interface %s: %s", errLinkByName, iface, err)
	}

	route := netlink.Route{
		Dst:       &destination,
		Gw:        gateway,
		LinkIndex: link.Attrs().Index,
		Table:     table,
	}
	if err := r.netLinker.RouteReplace(&route); err != nil {
		return fmt.Errorf("%w: for subnet %s at interface %s",
			err, destinationStr, iface)
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
		return fmt.Errorf("%w: for interface %s: %s", errLinkByName, iface, err)
	}

	route := netlink.Route{
		Dst:       &destination,
		Gw:        gateway,
		LinkIndex: link.Attrs().Index,
		Table:     table,
	}
	if err := r.netLinker.RouteDel(&route); err != nil {
		return fmt.Errorf("%w: for subnet %s at interface %s",
			err, destinationStr, iface)
	}

	return nil
}
