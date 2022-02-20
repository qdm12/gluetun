package routing

import (
	"fmt"
)

type Setuper interface {
	Setup() (err error)
}

func (r *Routing) Setup() (err error) {
	defaultInterfaceName, defaultGateway, err := r.DefaultRoute()
	if err != nil {
		return fmt.Errorf("cannot get default route: %w", err)
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
		return fmt.Errorf("cannot add routes for inbound traffic from default IP: %w", err)
	}

	r.stateMutex.RLock()
	outboundSubnets := r.outboundSubnets
	r.stateMutex.RUnlock()
	if err := r.setOutboundRoutes(outboundSubnets, defaultInterfaceName, defaultGateway); err != nil {
		return fmt.Errorf("cannot set outbound subnets routes: %w", err)
	}

	return nil
}

type TearDowner interface {
	TearDown() error
}

func (r *Routing) TearDown() error {
	defaultInterfaceName, defaultGateway, err := r.DefaultRoute()
	if err != nil {
		return fmt.Errorf("cannot get default route: %w", err)
	}

	err = r.unrouteInboundFromDefault(defaultGateway, defaultInterfaceName)
	if err != nil {
		return fmt.Errorf("cannot remove routes for inbound traffic from default IP: %w", err)
	}

	if err := r.setOutboundRoutes(nil, defaultInterfaceName, defaultGateway); err != nil {
		return fmt.Errorf("cannot set outbound subnets routes: %w", err)
	}

	return nil
}
