package routing

import (
	"fmt"
)

func (r *Routing) Setup() (err error) {
	defaultRoutes, err := r.DefaultRoutes()
	if err != nil {
		return fmt.Errorf("getting default routes: %w", err)
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

	err = r.routeInboundFromDefault(defaultRoutes)
	if err != nil {
		return fmt.Errorf("adding routes for inbound traffic from default IP: %w", err)
	}

	r.stateMutex.RLock()
	outboundSubnets := r.outboundSubnets
	r.stateMutex.RUnlock()
	if err := r.setOutboundRoutes(outboundSubnets, defaultRoutes); err != nil {
		return fmt.Errorf("setting outbound subnets routes: %w", err)
	}

	return nil
}

func (r *Routing) TearDown() error {
	defaultRoutes, err := r.DefaultRoutes()
	if err != nil {
		return fmt.Errorf("getting default route: %w", err)
	}

	err = r.unrouteInboundFromDefault(defaultRoutes)
	if err != nil {
		return fmt.Errorf("removing routes for inbound traffic from default IP: %w", err)
	}

	if err := r.setOutboundRoutes(nil, defaultRoutes); err != nil {
		return fmt.Errorf("setting outbound subnets routes: %w", err)
	}

	return nil
}
