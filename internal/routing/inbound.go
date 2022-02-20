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

func (r *Routing) routeInboundFromDefault(defaultGateway net.IP,
	defaultInterface string) (err error) {
	if err := r.addRuleInboundFromDefault(inboundTable); err != nil {
		return fmt.Errorf("cannot add rule: %w", err)
	}

	defaultDestination := net.IPNet{IP: net.IPv4(0, 0, 0, 0), Mask: net.IPv4Mask(0, 0, 0, 0)}
	if err := r.addRouteVia(defaultDestination, defaultGateway, defaultInterface, inboundTable); err != nil {
		return fmt.Errorf("cannot add route: %w", err)
	}

	return nil
}

func (r *Routing) unrouteInboundFromDefault(defaultGateway net.IP,
	defaultInterface string) (err error) {
	defaultDestination := net.IPNet{IP: net.IPv4(0, 0, 0, 0), Mask: net.IPv4Mask(0, 0, 0, 0)}
	if err := r.deleteRouteVia(defaultDestination, defaultGateway, defaultInterface, inboundTable); err != nil {
		return fmt.Errorf("cannot delete route: %w", err)
	}

	if err := r.delRuleInboundFromDefault(inboundTable); err != nil {
		return fmt.Errorf("cannot delete rule: %w", err)
	}

	return nil
}

func (r *Routing) addRuleInboundFromDefault(table int) (err error) {
	defaultIP, err := r.DefaultIP()
	if err != nil {
		return fmt.Errorf("cannot find default IP: %w", err)
	}

	defaultIPMasked32 := netlink.NewIPNet(defaultIP)
	ruleDstNet := (*net.IPNet)(nil)
	err = r.addIPRule(defaultIPMasked32, ruleDstNet, table, inboundPriority)
	if err != nil {
		return fmt.Errorf("cannot add rule: %w", err)
	}

	return nil
}

func (r *Routing) delRuleInboundFromDefault(table int) (err error) {
	defaultIP, err := r.DefaultIP()
	if err != nil {
		return fmt.Errorf("cannot find default IP: %w", err)
	}

	defaultIPMasked32 := netlink.NewIPNet(defaultIP)
	ruleDstNet := (*net.IPNet)(nil)
	err = r.deleteIPRule(defaultIPMasked32, ruleDstNet, table, inboundPriority)
	if err != nil {
		return fmt.Errorf("cannot delete rule: %w", err)
	}

	return nil
}
