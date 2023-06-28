package wireguard

import (
	"errors"
	"fmt"
	"net/netip"
	"strings"

	"github.com/qdm12/gluetun/internal/netlink"
)

func (w *Wireguard) addRoutes(link netlink.Link, destinations []netip.Prefix,
	firewallMark uint32,
) (err error) {
	for _, dst := range destinations {
		err = w.addRoute(link, dst, firewallMark)
		if err == nil {
			continue
		}

		if dst.Addr().Is6() && strings.Contains(err.Error(), "permission denied") {
			w.logger.Errorf("cannot add route for IPv6 due to a permission denial. "+
				"Ignoring and continuing execution; "+
				"Please report to https://github.com/qdm12/gluetun/issues/998 if you find a fix. "+
				"Full error string: %s", err)
			continue
		}
		return fmt.Errorf("adding route for destination %s: %w", dst, err)
	}
	return nil
}

var (
	ErrDefaultRouteNotFound = errors.New("default route not found")
)

func (w *Wireguard) addRoute(link netlink.Link, dst netip.Prefix,
	firewallMark uint32,
) (err error) {
	route := netlink.Route{
		LinkIndex: link.Index,
		Dst:       dst,
		Table:     int(firewallMark),
	}

	err = w.netlink.RouteAdd(route)
	if err != nil {
		return fmt.Errorf(
			"adding route for link %s, destination %s and table %d: %w",
			link.Name, dst, firewallMark, err)
	}

	vpnGatewayIP, err := w.routing.VPNLocalGatewayIP(link.Name)
	if err != nil {
		return fmt.Errorf("getting VPN gateway IP: %w", err)
	}

	routes, err := w.netlink.RouteList(netlink.FamilyV4)
	if err != nil {
		return fmt.Errorf("listing routes: %w", err)
	}

	var defaultRoute netlink.Route
	var defaultRouteFound bool
	for _, route = range routes {
		if !route.Dst.IsValid() || route.Dst.Addr().IsUnspecified() {
			defaultRoute = route
			defaultRouteFound = true
			break
		}
	}

	if !defaultRouteFound {
		return fmt.Errorf("%w: in %d routes", ErrDefaultRouteNotFound, len(routes))
	}

	// Equivalent replacement to:
	// ip route replace default via <vpn-gateway> dev tun0
	defaultRoute.Gw = vpnGatewayIP
	defaultRoute.LinkIndex = link.Index

	err = w.netlink.RouteReplace(defaultRoute)
	if err != nil {
		return fmt.Errorf("replacing default route: %w", err)
	}

	return err
}
