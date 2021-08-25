package routing

import (
	"bytes"
	"errors"
	"fmt"
	"net"

	"github.com/qdm12/gluetun/internal/netlink"
)

var (
	ErrVPNDestinationIPNotFound  = errors.New("VPN destination IP address not found")
	ErrVPNLocalGatewayIPNotFound = errors.New("VPN local gateway IP address not found")
)

type VPNDestinationIPGetter interface {
	VPNDestinationIP() (ip net.IP, err error)
}

func (r *Routing) VPNDestinationIP() (ip net.IP, err error) {
	routes, err := r.netLinker.RouteList(nil, netlink.FAMILY_ALL)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrRoutesList, err)
	}

	defaultLinkIndex := -1
	for _, route := range routes {
		if route.Dst == nil {
			defaultLinkIndex = route.LinkIndex
			break
		}
	}
	if defaultLinkIndex == -1 {
		return nil, fmt.Errorf("%w: in %d route(s)", ErrLinkDefaultNotFound, len(routes))
	}

	for _, route := range routes {
		if route.LinkIndex == defaultLinkIndex &&
			route.Dst != nil &&
			!IPIsPrivate(route.Dst.IP) &&
			bytes.Equal(route.Dst.Mask, net.IPMask{255, 255, 255, 255}) {
			return route.Dst.IP, nil
		}
	}
	return nil, fmt.Errorf("%w: in %d routes", ErrVPNDestinationIPNotFound, len(routes))
}

type VPNLocalGatewayIPGetter interface {
	VPNLocalGatewayIP(vpnIntf string) (ip net.IP, err error)
}

func (r *Routing) VPNLocalGatewayIP(vpnIntf string) (ip net.IP, err error) {
	routes, err := r.netLinker.RouteList(nil, netlink.FAMILY_ALL)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrRoutesList, err)
	}
	for _, route := range routes {
		link, err := r.netLinker.LinkByIndex(route.LinkIndex)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrLinkByIndex, err)
		}
		interfaceName := link.Attrs().Name
		if interfaceName == vpnIntf &&
			route.Dst != nil &&
			route.Dst.IP.Equal(net.IP{0, 0, 0, 0}) {
			return route.Gw, nil
		}
	}
	return nil, fmt.Errorf("%w: in %d routes", ErrVPNLocalGatewayIPNotFound, len(routes))
}
