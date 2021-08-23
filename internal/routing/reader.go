package routing

import (
	"bytes"
	"errors"
	"fmt"
	"net"

	"github.com/vishvananda/netlink"
)

type LocalNetwork struct {
	IPNet         *net.IPNet
	InterfaceName string
	IP            net.IP
}

var (
	ErrInterfaceIPNotFound       = errors.New("IP address not found for interface")
	ErrInterfaceListAddr         = errors.New("cannot list interface addresses")
	ErrInterfaceNotFound         = errors.New("network interface not found")
	ErrLinkByIndex               = errors.New("cannot obtain link by index")
	ErrLinkByName                = errors.New("cannot obtain link by name")
	ErrLinkDefaultNotFound       = errors.New("default link not found")
	ErrLinkList                  = errors.New("cannot list links")
	ErrLinkLocalNotFound         = errors.New("local link not found")
	ErrRouteDefaultNotFound      = errors.New("default route not found")
	ErrRoutesList                = errors.New("cannot list routes")
	ErrRulesList                 = errors.New("cannot list rules")
	ErrSubnetDefaultNotFound     = errors.New("default subnet not found")
	ErrSubnetLocalNotFound       = errors.New("local subnet not found")
	ErrVPNDestinationIPNotFound  = errors.New("VPN destination IP address not found")
	ErrVPNLocalGatewayIPNotFound = errors.New("VPN local gateway IP address not found")
)

type DefaultRouteGetter interface {
	DefaultRoute() (defaultInterface string, defaultGateway net.IP, err error)
}

func (r *Routing) DefaultRoute() (defaultInterface string, defaultGateway net.IP, err error) {
	routes, err := netlink.RouteList(nil, netlink.FAMILY_ALL)
	if err != nil {
		return "", nil, fmt.Errorf("%w: %s", ErrRoutesList, err)
	}
	for _, route := range routes {
		if route.Dst == nil {
			defaultGateway = route.Gw
			linkIndex := route.LinkIndex
			link, err := netlink.LinkByIndex(linkIndex)
			if err != nil {
				return "", nil, fmt.Errorf("%w: for default route at index %d: %s", ErrLinkByIndex, linkIndex, err)
			}
			attributes := link.Attrs()
			defaultInterface = attributes.Name
			r.logger.Info("default route found: interface " + defaultInterface +
				", gateway " + defaultGateway.String())
			return defaultInterface, defaultGateway, nil
		}
	}
	return "", nil, fmt.Errorf("%w: in %d route(s)", ErrRouteDefaultNotFound, len(routes))
}

type DefaultIPGetter interface {
	DefaultIP() (defaultIP net.IP, err error)
}

func (r *Routing) DefaultIP() (ip net.IP, err error) {
	routes, err := netlink.RouteList(nil, netlink.FAMILY_ALL)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrRoutesList, err)
	}

	defaultLinkName := ""
	for _, route := range routes {
		if route.Dst == nil {
			linkIndex := route.LinkIndex
			link, err := netlink.LinkByIndex(linkIndex)
			if err != nil {
				return nil, fmt.Errorf("%w: for default route at index %d: %s", ErrLinkByIndex, linkIndex, err)
			}
			defaultLinkName = link.Attrs().Name
		}
	}
	if defaultLinkName == "" {
		return nil, fmt.Errorf("%w: in %d route(s)", ErrLinkDefaultNotFound, len(routes))
	}

	return r.assignedIP(defaultLinkName)
}

type LocalSubnetGetter interface {
	LocalSubnet() (defaultSubnet net.IPNet, err error)
}

func (r *Routing) LocalSubnet() (defaultSubnet net.IPNet, err error) {
	routes, err := netlink.RouteList(nil, netlink.FAMILY_ALL)
	if err != nil {
		return defaultSubnet, fmt.Errorf("%w: %s", ErrRoutesList, err)
	}

	defaultLinkIndex := -1
	for _, route := range routes {
		if route.Dst == nil {
			defaultLinkIndex = route.LinkIndex
			break
		}
	}
	if defaultLinkIndex == -1 {
		return defaultSubnet, fmt.Errorf("%w: in %d route(s)", ErrLinkDefaultNotFound, len(routes))
	}

	for _, route := range routes {
		if route.Gw != nil || route.LinkIndex != defaultLinkIndex {
			continue
		}
		defaultSubnet = *route.Dst
		r.logger.Info("local subnet found: " + defaultSubnet.String())
		return defaultSubnet, nil
	}

	return defaultSubnet, fmt.Errorf("%w: in %d routes", ErrSubnetDefaultNotFound, len(routes))
}

type LocalNetworksGetter interface {
	LocalNetworks() (localNetworks []LocalNetwork, err error)
}

func (r *Routing) LocalNetworks() (localNetworks []LocalNetwork, err error) {
	links, err := netlink.LinkList()
	if err != nil {
		return localNetworks, fmt.Errorf("%w: %s", ErrLinkList, err)
	}

	localLinks := make(map[int]struct{})

	for _, link := range links {
		if link.Attrs().EncapType != "ether" {
			continue
		}

		localLinks[link.Attrs().Index] = struct{}{}
		r.logger.Info("local ethernet link found: " + link.Attrs().Name)
	}

	if len(localLinks) == 0 {
		return localNetworks, fmt.Errorf("%w: in %d links", ErrLinkLocalNotFound, len(links))
	}

	routes, err := netlink.RouteList(nil, netlink.FAMILY_V4)
	if err != nil {
		return localNetworks, fmt.Errorf("%w: %s", ErrRoutesList, err)
	}

	for _, route := range routes {
		if route.Gw != nil || route.Dst == nil {
			continue
		} else if _, ok := localLinks[route.LinkIndex]; !ok {
			continue
		}

		var localNet LocalNetwork

		localNet.IPNet = route.Dst
		r.logger.Info("local ipnet found: " + localNet.IPNet.String())

		link, err := netlink.LinkByIndex(route.LinkIndex)
		if err != nil {
			return localNetworks, fmt.Errorf("%w: at index %d: %s", ErrLinkByIndex, route.LinkIndex, err)
		}

		localNet.InterfaceName = link.Attrs().Name

		ip, err := r.assignedIP(localNet.InterfaceName)
		if err != nil {
			return localNetworks, err
		}

		localNet.IP = ip

		localNetworks = append(localNetworks, localNet)
	}

	if len(localNetworks) == 0 {
		return localNetworks, fmt.Errorf("%w: in %d routes", ErrSubnetLocalNotFound, len(routes))
	}

	return localNetworks, nil
}

func (r *Routing) assignedIP(interfaceName string) (ip net.IP, err error) {
	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return nil, fmt.Errorf("%w: %s: %s", ErrInterfaceNotFound, interfaceName, err)
	}
	addresses, err := iface.Addrs()
	if err != nil {
		return nil, fmt.Errorf("%w: %s: %s", ErrInterfaceListAddr, interfaceName, err)
	}
	for _, address := range addresses {
		switch value := address.(type) {
		case *net.IPAddr:
			return value.IP, nil
		case *net.IPNet:
			return value.IP, nil
		}
	}
	return nil, fmt.Errorf("%w: interface %s in %d addresses",
		ErrInterfaceIPNotFound, interfaceName, len(addresses))
}

type VPNDestinationIPGetter interface {
	VPNDestinationIP() (ip net.IP, err error)
}

func (r *Routing) VPNDestinationIP() (ip net.IP, err error) {
	routes, err := netlink.RouteList(nil, netlink.FAMILY_ALL)
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
	routes, err := netlink.RouteList(nil, netlink.FAMILY_ALL)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrRoutesList, err)
	}
	for _, route := range routes {
		link, err := netlink.LinkByIndex(route.LinkIndex)
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

func IPIsPrivate(ip net.IP) bool {
	return ip.IsPrivate() || ip.IsLoopback() ||
		ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast()
}
