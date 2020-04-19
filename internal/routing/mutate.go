package routing

import (
	"context"
	"net"

	"fmt"
)

func (r *routing) AddRoutesVia(ctx context.Context, subnets []net.IPNet, defaultGateway net.IP, defaultInterface string) error {
	for _, subnet := range subnets {
		exists, err := r.routeExists(subnet)
		if err != nil {
			return err
		} else if exists { // thanks to @npawelek https://github.com/npawelek
			if err := r.removeRoute(ctx, subnet); err != nil {
				return err
			}
		}
		r.logger.Info("adding %s as route via %s", subnet.String(), defaultInterface)
		output, err := r.commander.Run(ctx, "ip", "route", "add", subnet.String(), "via", defaultGateway.String(), "dev", defaultInterface)
		if err != nil {
			return fmt.Errorf("cannot add route for %s via %s %s %s: %s: %w", subnet.String(), defaultGateway.String(), "dev", defaultInterface, output, err)
		}
	}
	return nil
}

func (r *routing) removeRoute(ctx context.Context, subnet net.IPNet) (err error) {
	output, err := r.commander.Run(ctx, "ip", "route", "del", subnet.String())
	if err != nil {
		return fmt.Errorf("cannot delete route for %s: %s: %w", subnet.String(), output, err)
	}
	return nil
}
