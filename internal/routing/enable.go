package routing

import (
	"fmt"
	"net"
)

var (
	ErrSetup    = fmt.Errorf("cannot setup routing")
	ErrTeardown = fmt.Errorf("cannot teardown routing")
)

const (
	table    = 200
	priority = 100
)

func (r *routing) Setup() (err error) {
	defaultIP, err := r.defaultIP()
	if err != nil {
		return fmt.Errorf("%s: %w", ErrSetup, err)
	}
	defaultInterfaceName, defaultGateway, err := r.DefaultRoute()
	if err != nil {
		return fmt.Errorf("%s: %w", ErrSetup, err)
	}

	defer func() {
		if err == nil {
			return
		}
		if err := r.TearDown(); err != nil {
			r.logger.Error(err)
		}
	}()
	if err := r.addIPRule(defaultIP, table, priority); err != nil {
		return fmt.Errorf("%s: %w", ErrSetup, err)
	}
	defaultDestination := net.IPNet{IP: net.IPv4(0, 0, 0, 0), Mask: net.IPv4Mask(0, 0, 0, 0)}
	if err := r.addRouteVia(defaultDestination, defaultGateway, defaultInterfaceName, table); err != nil {
		return fmt.Errorf("%s: %w", ErrSetup, err)
	}

	r.stateMutex.RLock()
	outboundSubnets := r.outboundSubnets
	r.stateMutex.RUnlock()
	if err := r.setOutboundRoutes(outboundSubnets, defaultInterfaceName, defaultGateway); err != nil {
		return fmt.Errorf("%s: %w", ErrSetup, err)
	}

	return nil
}

func (r *routing) TearDown() error {
	defaultIP, err := r.defaultIP()
	if err != nil {
		return fmt.Errorf("%s: %w", ErrTeardown, err)
	}
	defaultInterfaceName, defaultGateway, err := r.DefaultRoute()
	if err != nil {
		return fmt.Errorf("%s: %w", ErrTeardown, err)
	}

	defaultNet := net.IPNet{IP: net.IPv4(0, 0, 0, 0), Mask: net.IPv4Mask(0, 0, 0, 0)}
	if err := r.deleteRouteVia(defaultNet, defaultGateway, defaultInterfaceName, table); err != nil {
		return fmt.Errorf("%s: %w", ErrTeardown, err)
	}
	if err := r.deleteIPRule(defaultIP, table, priority); err != nil {
		return fmt.Errorf("%s: %w", ErrTeardown, err)
	}

	if err := r.setOutboundRoutes(nil, defaultInterfaceName, defaultGateway); err != nil {
		return fmt.Errorf("%s: %w", ErrSetup, err)
	}

	return nil
}
