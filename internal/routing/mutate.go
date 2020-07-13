package routing

import (
	"context"
	"net"

	"fmt"
)

func (r *routing) AddRouteVia(ctx context.Context, subnet net.IPNet, defaultGateway net.IP, defaultInterface string) error {
	subnetStr := subnet.String()
	r.logger.Info("adding %s as route via %s %s", subnetStr, defaultGateway, defaultInterface)
	exists, err := r.routeExists(subnet)
	if err != nil {
		return err
	} else if exists {
		return nil
	}
	if r.debug {
		fmt.Printf("ip route add %s via %s dev %s\n", subnetStr, defaultGateway, defaultInterface)
	}
	output, err := r.commander.Run(ctx, "ip", "route", "add", subnetStr, "via", defaultGateway.String(), "dev", defaultInterface)
	if err != nil {
		return fmt.Errorf("cannot add route for %s via %s %s %s: %s: %w", subnetStr, defaultGateway, "dev", defaultInterface, output, err)
	}
	return nil
}

func (r *routing) DeleteRouteVia(ctx context.Context, subnet net.IPNet) (err error) {
	subnetStr := subnet.String()
	r.logger.Info("deleting route for %s", subnetStr)
	exists, err := r.routeExists(subnet)
	if err != nil {
		return err
	} else if !exists { // thanks to @npawelek https://github.com/npawelek
		return nil
	}
	if r.debug {
		fmt.Printf("ip route del %s\n", subnetStr)
	}
	output, err := r.commander.Run(ctx, "ip", "route", "del", subnetStr)
	if err != nil {
		return fmt.Errorf("cannot delete route for %s: %s: %w", subnetStr, output, err)
	}
	return nil
}
