package routing

import (
	"net"

	"fmt"
)

func (r *routing) AddRoutesVia(subnets []net.IPNet, defaultGateway net.IP, defaultInterface string) error {
	for _, subnet := range subnets {
		exists, err := r.routeExists(subnet)
		if err != nil {
			return fmt.Errorf("cannot check route for subnet %s: %w", subnet, err)
		} else if exists { // thanks to @npawelek https://github.com/npawelek
			continue
		}
		r.logger.Info("adding %s as route via %s", subnet.String(), defaultInterface)
		output, err := r.commander.Run("ip", "route", "add", subnet.String(), "via", defaultGateway.String(), "dev", defaultInterface)
		if err != nil {
			return fmt.Errorf("cannot add route for %s via %s %s %s: %s: %w", subnet.String(), defaultGateway.String(), "dev", defaultInterface, output, err)
		}
	}
	return nil
}
