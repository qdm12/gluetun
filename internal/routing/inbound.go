package routing

import (
	"fmt"
	"net/netip"

	"github.com/qdm12/gluetun/internal/netlink"
)

const (
	inboundTable    = 200
	inboundPriority = 100
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
		if defaultRoute.Family == netlink.FAMILY_V6 {
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
		if defaultRoute.Family == netlink.FAMILY_V6 {
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

func (r *Routing) addRuleInboundFromDefault(table int, defaultRoutes []DefaultRoute) (err error) {
	for _, defaultRoute := range defaultRoutes {
		assignedIP := netIPToNetipAddress(defaultRoute.AssignedIP)
		bits := 32
		if assignedIP.Is6() {
			bits = 128
		}
		r.logger.Debug(fmt.Sprintf("ASSIGNED IP IS %#v -> %s, bits %d",
			defaultRoute.AssignedIP, assignedIP, bits))
		defaultIPMasked := netip.PrefixFrom(assignedIP, bits)
		ruleDstNet := (*netip.Prefix)(nil)
		err = r.addIPRule(&defaultIPMasked, ruleDstNet, table, inboundPriority)
		if err != nil {
			return fmt.Errorf("adding rule for default route %s: %w", defaultRoute, err)
		}
	}

	return nil
}

func (r *Routing) delRuleInboundFromDefault(table int, defaultRoutes []DefaultRoute) (err error) {
	for _, defaultRoute := range defaultRoutes {
		assignedIP := netIPToNetipAddress(defaultRoute.AssignedIP)
		bits := 32
		if assignedIP.Is6() {
			bits = 128
		}
		defaultIPMasked := netip.PrefixFrom(assignedIP, bits)
		ruleDstNet := (*netip.Prefix)(nil)
		err = r.deleteIPRule(&defaultIPMasked, ruleDstNet, table, inboundPriority)
		if err != nil {
			return fmt.Errorf("deleting rule for default route %s: %w", defaultRoute, err)
		}
	}

	return nil
}
