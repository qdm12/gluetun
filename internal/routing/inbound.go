package routing

import (
	"errors"
	"fmt"
	"net"

	"github.com/qdm12/gluetun/internal/netlink"
)

const (
	inboundTable    = 200
	inboundPriority = 100
)

var (
	errDefaultIP = errors.New("cannot get default IP address")
)

func (r *Routing) routeInboundFromDefault(defaultGateway net.IP,
	defaultInterface string) (err error) {
	if err := r.addRuleInboundFromDefault(inboundTable); err != nil {
		return fmt.Errorf("%w: %s", errRuleAdd, err)
	}

	defaultDestination := net.IPNet{IP: net.IPv4(0, 0, 0, 0), Mask: net.IPv4Mask(0, 0, 0, 0)}
	if err := r.addRouteVia(defaultDestination, defaultGateway, defaultInterface, inboundTable); err != nil {
		return fmt.Errorf("%w: %s", errRouteAdd, err)
	}

	return nil
}

func (r *Routing) unrouteInboundFromDefault(defaultGateway net.IP,
	defaultInterface string) (err error) {
	defaultDestination := net.IPNet{IP: net.IPv4(0, 0, 0, 0), Mask: net.IPv4Mask(0, 0, 0, 0)}
	if err := r.deleteRouteVia(defaultDestination, defaultGateway, defaultInterface, inboundTable); err != nil {
		return fmt.Errorf("%w: %s", errRouteDelete, err)
	}

	if err := r.delRuleInboundFromDefault(inboundTable); err != nil {
		return fmt.Errorf("%w: %s", errRuleDelete, err)
	}

	return nil
}

func (r *Routing) addRuleInboundFromDefault(table int) (err error) {
	defaultIP, err := r.DefaultIP()
	if err != nil {
		return fmt.Errorf("%w: %s", errDefaultIP, err)
	}

	defaultIPMasked32 := netlink.NewIPNet(defaultIP)
	ruleDstNet := (*net.IPNet)(nil)
	err = r.addIPRule(defaultIPMasked32, ruleDstNet, table, inboundPriority)
	if err != nil {
		return fmt.Errorf("%w: %s", errRuleAdd, err)
	}

	return nil
}

func (r *Routing) delRuleInboundFromDefault(table int) (err error) {
	defaultIP, err := r.DefaultIP()
	if err != nil {
		return fmt.Errorf("%w: %s", errDefaultIP, err)
	}

	defaultIPMasked32 := netlink.NewIPNet(defaultIP)
	ruleDstNet := (*net.IPNet)(nil)
	err = r.deleteIPRule(defaultIPMasked32, ruleDstNet, table, inboundPriority)
	if err != nil {
		return fmt.Errorf("%w: %s", errRuleDelete, err)
	}

	return nil
}
