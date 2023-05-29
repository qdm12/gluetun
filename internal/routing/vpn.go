package routing

import (
	"errors"
	"fmt"
	"net/netip"

	"github.com/qdm12/gluetun/internal/netlink"
)

var (
	ErrVPNDestinationIPNotFound  = errors.New("VPN destination IP address not found")
	ErrVPNLocalGatewayIPNotFound = errors.New("VPN local gateway IP address not found")
)

func (r *Routing) VPNDestinationIP() (ip netip.Addr, err error) {
	routes, err := r.netLinker.RouteList(nil, netlink.FamilyAll)
	if err != nil {
		return ip, fmt.Errorf("listing routes: %w", err)
	}

	defaultLinkIndex := -1
	for _, route := range routes {
		if !route.Dst.IsValid() {
			defaultLinkIndex = route.LinkIndex
			break
		}
	}
	if defaultLinkIndex == -1 {
		return ip, fmt.Errorf("%w: in %d route(s)", ErrLinkDefaultNotFound, len(routes))
	}

	for _, route := range routes {
		if route.LinkIndex == defaultLinkIndex &&
			route.Dst.IsValid() &&
			!IPIsPrivate(route.Dst.Addr()) &&
			route.Dst.IsSingleIP() {
			return route.Dst.Addr(), nil
		}
	}
	return ip, fmt.Errorf("%w: in %d routes", ErrVPNDestinationIPNotFound, len(routes))
}

func (r *Routing) VPNLocalGatewayIP(vpnIntf string) (ip netip.Addr, err error) {
	routes, err := r.netLinker.RouteList(nil, netlink.FamilyAll)
	if err != nil {
		return ip, fmt.Errorf("listing routes: %w", err)
	}
	for _, route := range routes {
		link, err := r.netLinker.LinkByIndex(route.LinkIndex)
		if err != nil {
			return ip, fmt.Errorf("finding link at index %d: %w", route.LinkIndex, err)
		}
		interfaceName := link.Name
		if interfaceName == vpnIntf &&
			route.Dst.IsValid() &&
			route.Dst.Addr().IsUnspecified() {
			return route.Gw, nil
		}
	}
	return ip, fmt.Errorf("%w: in %d routes", ErrVPNLocalGatewayIPNotFound, len(routes))
}
