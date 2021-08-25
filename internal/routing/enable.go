package routing

import (
	"errors"
	"fmt"
)

var (
	ErrDefaultRoute          = errors.New("cannot get default route")
	ErrAddInboundFromDefault = errors.New("cannot add routes for inbound traffic from default IP")
	ErrDelInboundFromDefault = errors.New("cannot remove routes for inbound traffic from default IP")
	ErrSubnetsOutboundSet    = errors.New("cannot set outbound subnets routes")
)

type Setuper interface {
	Setup() (err error)
}

func (r *Routing) Setup() (err error) {
	defaultInterfaceName, defaultGateway, err := r.DefaultRoute()
	if err != nil {
		return fmt.Errorf("%w: %s", ErrDefaultRoute, err)
	}

	touched := false
	defer func() {
		if err != nil && touched {
			if tearDownErr := r.TearDown(); tearDownErr != nil {
				r.logger.Error("cannot reverse routing changes: " + tearDownErr.Error())
			}
		}
	}()

	touched = true

	err = r.routeInboundFromDefault(defaultGateway, defaultInterfaceName)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrAddInboundFromDefault, err)
	}

	r.stateMutex.RLock()
	outboundSubnets := r.outboundSubnets
	r.stateMutex.RUnlock()
	if err := r.setOutboundRoutes(outboundSubnets, defaultInterfaceName, defaultGateway); err != nil {
		return fmt.Errorf("%w: %s", ErrSubnetsOutboundSet, err)
	}

	return nil
}

type TearDowner interface {
	TearDown() error
}

func (r *Routing) TearDown() error {
	defaultInterfaceName, defaultGateway, err := r.DefaultRoute()
	if err != nil {
		return fmt.Errorf("%w: %s", ErrDefaultRoute, err)
	}

	err = r.unrouteInboundFromDefault(defaultGateway, defaultInterfaceName)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrDelInboundFromDefault, err)
	}

	if err := r.setOutboundRoutes(nil, defaultInterfaceName, defaultGateway); err != nil {
		return fmt.Errorf("%w: %s", ErrSubnetsOutboundSet, err)
	}

	return nil
}
