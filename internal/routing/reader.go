package routing

import (
	"bytes"
	"fmt"
	"net"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/vishvananda/netlink"
)

func (r *routing) DefaultRoute() (defaultInterface string, defaultGateway net.IP, err error) {
	routes, err := netlink.RouteList(nil, netlink.FAMILY_ALL)
	if err != nil {
		return "", nil, fmt.Errorf("cannot list routes: %w", err)
	}
	for _, route := range routes {
		if route.Dst == nil {
			defaultGateway = route.Gw
			linkIndex := route.LinkIndex
			link, err := netlink.LinkByIndex(linkIndex)
			if err != nil {
				return "", nil, fmt.Errorf("cannot obtain link with index %d for default route: %w", linkIndex, err)
			}
			attributes := link.Attrs()
			defaultInterface = attributes.Name
			r.logger.Info("default route found: interface %s, gateway %s", defaultInterface, defaultGateway.String())
			return defaultInterface, defaultGateway, nil
		}
	}
	return "", nil, fmt.Errorf("cannot find default route in %d routes", len(routes))
}

func (r *routing) LocalSubnet() (defaultSubnet net.IPNet, err error) {
	routes, err := netlink.RouteList(nil, netlink.FAMILY_ALL)
	if err != nil {
		return defaultSubnet, fmt.Errorf("cannot find local subnet: %w", err)
	}

	defaultLinkIndex := -1
	for _, route := range routes {
		if route.Dst == nil {
			defaultLinkIndex = route.LinkIndex
			break
		}
	}
	if defaultLinkIndex == -1 {
		return defaultSubnet, fmt.Errorf("cannot find local subnet: cannot find default link")
	}

	for _, route := range routes {
		if route.Gw != nil || route.LinkIndex != defaultLinkIndex {
			continue
		}
		defaultSubnet = *route.Dst
		r.logger.Info("local subnet found: %s", defaultSubnet.String())
		return defaultSubnet, nil
	}

	return defaultSubnet, fmt.Errorf("cannot find default subnet in %d routes", len(routes))
}

func (r *routing) routeExists(subnet net.IPNet) (exists bool, err error) {
	routes, err := netlink.RouteList(nil, netlink.FAMILY_ALL)
	if err != nil {
		return false, fmt.Errorf("cannot list routes: %w", err)
	}
	for _, route := range routes {
		if route.Dst != nil && route.Dst.String() == subnet.String() {
			return true, nil
		}
	}
	return false, nil
}

func (r *routing) VPNDestinationIP() (ip net.IP, err error) {
	routes, err := netlink.RouteList(nil, netlink.FAMILY_ALL)
	if err != nil {
		return nil, fmt.Errorf("cannot find VPN destination IP: %w", err)
	}

	defaultLinkIndex := -1
	for _, route := range routes {
		if route.Dst == nil {
			defaultLinkIndex = route.LinkIndex
			break
		}
	}
	if defaultLinkIndex == -1 {
		return nil, fmt.Errorf("cannot find VPN destination IP: cannot find default link")
	}

	for _, route := range routes {
		if route.LinkIndex == defaultLinkIndex &&
			route.Dst != nil &&
			!ipIsPrivate(route.Dst.IP) &&
			bytes.Equal(route.Dst.Mask, net.IPMask{255, 255, 255, 255}) {
			return route.Dst.IP, nil
		}
	}
	return nil, fmt.Errorf("cannot find VPN destination IP address from ip routes")
}

func (r *routing) VPNLocalGatewayIP() (ip net.IP, err error) {
	routes, err := netlink.RouteList(nil, netlink.FAMILY_ALL)
	if err != nil {
		return nil, fmt.Errorf("cannot find VPN local gateway IP: %w", err)
	}
	for _, route := range routes {
		link, err := netlink.LinkByIndex(route.LinkIndex)
		if err != nil {
			return nil, fmt.Errorf("cannot find VPN local gateway IP: %w", err)
		}
		interfaceName := link.Attrs().Name
		if interfaceName == string(constants.TUN) &&
			route.Dst != nil &&
			route.Dst.IP.Equal(net.IP{0, 0, 0, 0}) {
			return route.Gw, nil
		}
	}
	return nil, fmt.Errorf("cannot find VPN local gateway IP address from ip routes")
}

func ipIsPrivate(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}
	privateCIDRBlocks := [8]string{
		"127.0.0.0/8",    // localhost
		"10.0.0.0/8",     // 24-bit block
		"172.16.0.0/12",  // 20-bit block
		"192.168.0.0/16", // 16-bit block
		"169.254.0.0/16", // link local address
		"::1/128",        // localhost IPv6
		"fc00::/7",       // unique local address IPv6
		"fe80::/10",      // link local address IPv6
	}
	for i := range privateCIDRBlocks {
		_, CIDR, _ := net.ParseCIDR(privateCIDRBlocks[i])
		if CIDR.Contains(ip) {
			return true
		}
	}
	return false
}
