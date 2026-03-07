package routing

import (
	"fmt"
	"net/netip"

	"github.com/qdm12/gluetun/internal/netlink"
)

const (
	inboundTable    uint32 = 200
	inboundPriority uint32 = 100
)

func (r *Routing) routeInboundFromDefault(defaultRoutes []DefaultRoute) (err error) {
	if err := r.addRuleInboundFromDefault(inboundTable, defaultRoutes); err != nil {
		return fmt.Errorf("adding rule: %w", err)
	}

	const bits = 0
	defaultDestinationIPv4 := netip.PrefixFrom(netip.AddrFrom4([4]byte{}), bits)
	defaultDestinationIPv6 := netip.PrefixFrom(netip.AddrFrom16([16]byte{}), bits)

	for _, defaultRoute := range defaultRoutes {
		defaultDestination := defaultDestinationIPv4
		if defaultRoute.Family == netlink.FamilyV6 {
			defaultDestination = defaultDestinationIPv6
		}

		err := r.addRouteVia(defaultDestination, defaultRoute.Gateway, defaultRoute.NetInterface, inboundTable)
		if err != nil {
			return fmt.Errorf("adding route: %w", err)
		}
	}

	return nil
}

func (r *Routing) unrouteInboundFromDefault(defaultRoutes []DefaultRoute) (err error) {
	const bits = 0
	defaultDestinationIPv4 := netip.PrefixFrom(netip.AddrFrom4([4]byte{}), bits)
	defaultDestinationIPv6 := netip.PrefixFrom(netip.AddrFrom16([16]byte{}), bits)

	for _, defaultRoute := range defaultRoutes {
		defaultDestination := defaultDestinationIPv4
		if defaultRoute.Family == netlink.FamilyV6 {
			defaultDestination = defaultDestinationIPv6
		}

		err := r.deleteRouteVia(defaultDestination, defaultRoute.Gateway, defaultRoute.NetInterface, inboundTable)
		if err != nil {
			return fmt.Errorf("deleting route: %w", err)
		}
	}

	if err := r.delRuleInboundFromDefault(inboundTable, defaultRoutes); err != nil {
		return fmt.Errorf("deleting rule: %w", err)
	}

	return nil
}

func (r *Routing) addRuleInboundFromDefault(table uint32, defaultRoutes []DefaultRoute) (err error) {
	for _, defaultRoute := range defaultRoutes {
		assignedIP := defaultRoute.AssignedIP
		bits := 32
		if assignedIP.Is6() {
			bits = 128
		}
		defaultIPMasked := netip.PrefixFrom(assignedIP, bits)
		ruleDstNet := netip.Prefix{}
		err = r.addIPRule(defaultIPMasked, ruleDstNet, table, inboundPriority)
		if err != nil {
			return fmt.Errorf("adding rule for default route %s: %w", defaultRoute, err)
		}
	}

	return nil
}

func (r *Routing) delRuleInboundFromDefault(table uint32, defaultRoutes []DefaultRoute) (err error) {
	for _, defaultRoute := range defaultRoutes {
		assignedIP := defaultRoute.AssignedIP
		bits := 32
		if assignedIP.Is6() {
			bits = 128
		}
		defaultIPMasked := netip.PrefixFrom(assignedIP, bits)
		ruleDstNet := netip.Prefix{}
		err = r.deleteIPRule(defaultIPMasked, ruleDstNet, table, inboundPriority)
		if err != nil {
			return fmt.Errorf("deleting rule for default route %s: %w", defaultRoute, err)
		}
	}

	return nil
}
