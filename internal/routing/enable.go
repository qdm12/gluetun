package routing

import (
	"errors"
	"fmt"
	"net"
)

const (
	table    = 200
	priority = 100
)

var (
	ErrDefaultIP          = errors.New("cannot get default IP address")
	ErrDefaultRoute       = errors.New("cannot get default route")
	ErrIPRuleAdd          = errors.New("cannot add IP rule")
	ErrIPRuleDelete       = errors.New("cannot delete IP rule")
	ErrRouteAdd           = errors.New("cannot add route")
	ErrSubnetsOutboundSet = errors.New("cannot set outbound subnets routes")
)

func (r *routing) Setup() (err error) {
	defaultIP, err := r.DefaultIP()
	if err != nil {
		return fmt.Errorf("%w: %s", ErrDefaultIP, err)
	}
	defaultInterfaceName, defaultGateway, err := r.DefaultRoute()
	if err != nil {
		return fmt.Errorf("%w: %s", ErrDefaultRoute, err)
	}

	defer func() {
		if err == nil {
			return
		}
		if err := r.TearDown(); err != nil {
			r.logger.Error("cannot reverse routing changes: " + err.Error())
		}
	}()
	if err := r.addIPRule(defaultIP, table, priority); err != nil {
		return fmt.Errorf("%w: %s", ErrIPRuleAdd, err)
	}
	defaultDestination := net.IPNet{IP: net.IPv4(0, 0, 0, 0), Mask: net.IPv4Mask(0, 0, 0, 0)}
	if err := r.addRouteVia(defaultDestination, defaultGateway, defaultInterfaceName, table); err != nil {
		return fmt.Errorf("%w: %s", ErrRouteAdd, err)
	}

	r.stateMutex.RLock()
	outboundSubnets := r.outboundSubnets
	r.stateMutex.RUnlock()
	if err := r.setOutboundRoutes(outboundSubnets, defaultInterfaceName, defaultGateway); err != nil {
		return fmt.Errorf("%w: %s", ErrSubnetsOutboundSet, err)
	}

	return nil
}

func (r *routing) TearDown() error {
	defaultIP, err := r.DefaultIP()
	if err != nil {
		return fmt.Errorf("%w: %s", ErrDefaultIP, err)
	}
	defaultInterfaceName, defaultGateway, err := r.DefaultRoute()
	if err != nil {
		return fmt.Errorf("%w: %s", ErrDefaultRoute, err)
	}

	defaultNet := net.IPNet{IP: net.IPv4(0, 0, 0, 0), Mask: net.IPv4Mask(0, 0, 0, 0)}
	if err := r.deleteRouteVia(defaultNet, defaultGateway, defaultInterfaceName, table); err != nil {
		return fmt.Errorf("%w: %s", ErrRouteDelete, err)
	}
	if err := r.deleteIPRule(defaultIP, table, priority); err != nil {
		return fmt.Errorf("%w: %s", ErrIPRuleDelete, err)
	}

	if err := r.setOutboundRoutes(nil, defaultInterfaceName, defaultGateway); err != nil {
		return fmt.Errorf("%w: %s", ErrSubnetsOutboundSet, err)
	}

	return nil
}
