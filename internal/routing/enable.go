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
	if err := r.addRouteVia(net.IPNet{}, defaultGateway, defaultInterfaceName, table); err != nil {
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

	if err := r.deleteRouteVia(net.IPNet{}, defaultGateway, defaultInterfaceName, table); err != nil {
		return fmt.Errorf("%s: %w", ErrTeardown, err)
	}
	if err := r.deleteIPRule(defaultIP, table, priority); err != nil {
		return fmt.Errorf("%s: %w", ErrTeardown, err)
	}
	return nil
}
