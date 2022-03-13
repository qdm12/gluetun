package routing

import (
	"fmt"
	"net"

	"github.com/qdm12/gluetun/internal/netlink"
)

const (
	inboundTable    = 200
	inboundPriority = 100
)

func (r *Routing) routeInboundFromDefault(defaultRoutes []DefaultRoute) (err error) {
	if err := r.addRuleInboundFromDefault(inboundTable, defaultRoutes); err != nil {
		return fmt.Errorf("cannot add rule: %w", err)
	}

	defaultDestinationIPv4 := net.IPNet{IP: net.IPv4(0, 0, 0, 0), Mask: net.IPv4Mask(0, 0, 0, 0)}
	defaultDestinationIPv6 := net.IPNet{IP: net.IPv6zero, Mask: net.IPMask(net.IPv6zero)}

	for _, defaultRoute := range defaultRoutes {
		defaultDestination := defaultDestinationIPv4
		if defaultRoute.Family == netlink.FAMILY_V6 {
			defaultDestination = defaultDestinationIPv6
		}

		err := r.addRouteVia(defaultDestination, defaultRoute.Gateway, defaultRoute.NetInterface, inboundTable)
		if err != nil {
			return fmt.Errorf("cannot add route: %w", err)
		}
	}

	return nil
}

func (r *Routing) unrouteInboundFromDefault(defaultRoutes []DefaultRoute) (err error) {
	defaultDestinationIPv4 := net.IPNet{IP: net.IPv4(0, 0, 0, 0), Mask: net.IPv4Mask(0, 0, 0, 0)}
	defaultDestinationIPv6 := net.IPNet{IP: net.IPv6zero, Mask: net.IPMask(net.IPv6zero)}

	for _, defaultRoute := range defaultRoutes {
		defaultDestination := defaultDestinationIPv4
		if defaultRoute.Family == netlink.FAMILY_V6 {
			defaultDestination = defaultDestinationIPv6
		}

		err := r.deleteRouteVia(defaultDestination, defaultRoute.Gateway, defaultRoute.NetInterface, inboundTable)
		if err != nil {
			return fmt.Errorf("cannot delete route: %w", err)
		}
	}

	if err := r.delRuleInboundFromDefault(inboundTable, defaultRoutes); err != nil {
		return fmt.Errorf("cannot delete rule: %w", err)
	}

	return nil
}

func (r *Routing) addRuleInboundFromDefault(table int, defaultRoutes []DefaultRoute) (err error) {
	for _, defaultRoute := range defaultRoutes {
		defaultIPMasked32 := netlink.NewIPNet(defaultRoute.AssignedIP)
		ruleDstNet := (*net.IPNet)(nil)
		err = r.addIPRule(defaultIPMasked32, ruleDstNet, table, inboundPriority)
		if err != nil {
			return fmt.Errorf("cannot add rule for default route %s: %w", defaultRoute, err)
		}
	}

	return nil
}

func (r *Routing) delRuleInboundFromDefault(table int, defaultRoutes []DefaultRoute) (err error) {
	for _, defaultRoute := range defaultRoutes {
		defaultIPMasked32 := netlink.NewIPNet(defaultRoute.AssignedIP)
		ruleDstNet := (*net.IPNet)(nil)
		err = r.deleteIPRule(defaultIPMasked32, ruleDstNet, table, inboundPriority)
		if err != nil {
			return fmt.Errorf("cannot delete rule for default route %s: %w", defaultRoute, err)
		}
	}

	return nil
}
