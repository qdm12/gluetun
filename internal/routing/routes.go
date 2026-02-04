package routing

import (
	"fmt"
	"net/netip"
	"strconv"

	"github.com/qdm12/gluetun/internal/netlink"
)

func (r *Routing) addRouteVia(destination netip.Prefix, gateway netip.Addr,
	iface string, table uint32,
) error {
	destinationStr := destination.String()
	r.logger.Info("adding route for " + destinationStr)
	r.logger.Debug("ip route replace " + destinationStr +
		" via " + gateway.String() +
		" dev " + iface +
		" table " + strconv.Itoa(int(table)))

	link, err := r.netLinker.LinkByName(iface)
	if err != nil {
		return fmt.Errorf("finding link for interface %s: %w", iface, err)
	}

	family := netlink.FamilyV4
	if destination.Addr().Is6() {
		family = netlink.FamilyV6
	}
	route := netlink.Route{
		Dst:       destination,
		Gw:        gateway,
		LinkIndex: link.Index,
		Family:    family,
		Table:     table,
		Type:      netlink.RouteTypeUnicast,
		Scope:     netlink.ScopeUniverse,
		Proto:     netlink.ProtoStatic,
	}
	if err := r.netLinker.RouteReplace(route); err != nil {
		return fmt.Errorf("replacing route for subnet %s at interface %s: %w",
			destinationStr, iface, err)
	}

	return nil
}

func (r *Routing) deleteRouteVia(destination netip.Prefix, gateway netip.Addr,
	iface string, table uint32,
) (err error) {
	destinationStr := destination.String()
	r.logger.Info("deleting route for " + destinationStr)
	r.logger.Debug("ip route delete " + destinationStr +
		" via " + gateway.String() +
		" dev " + iface +
		" table " + strconv.Itoa(int(table)))

	link, err := r.netLinker.LinkByName(iface)
	if err != nil {
		return fmt.Errorf("finding link for interface %s: %w", iface, err)
	}

	family := netlink.FamilyV4
	if destination.Addr().Is6() {
		family = netlink.FamilyV6
	}
	route := netlink.Route{
		Dst:       destination,
		Gw:        gateway,
		LinkIndex: link.Index,
		Family:    family,
		Table:     table,
	}
	if err := r.netLinker.RouteDel(route); err != nil {
		return fmt.Errorf("deleting route: for subnet %s at interface %s: %w",
			destinationStr, iface, err)
	}

	return nil
}
