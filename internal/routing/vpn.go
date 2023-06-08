package routing

import (
	"errors"
	"fmt"
	"net/netip"

	"github.com/qdm12/gluetun/internal/netlink"
)

var (
	ErrVPNLocalGatewayIPNotFound = errors.New("VPN local gateway IP address not found")
)

func (r *Routing) VPNLocalGatewayIP(vpnIntf string) (ip netip.Addr, err error) {
	routes, err := r.netLinker.RouteList(netlink.FamilyAll)
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
