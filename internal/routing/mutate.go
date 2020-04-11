package routing

import (
	"net"

	"fmt"
)

func (r *routing) AddRoutesVia(subnets []net.IPNet, defaultGateway net.IP, defaultInterface string) error {
	for _, subnet := range subnets {
		subnetStr := subnet.String()
		output, err := r.commander.Run("ip", "route", "show", subnetStr)
		if err != nil {
			return fmt.Errorf("cannot read route %s: %s: %w", subnetStr, output, err)
		} else if len(output) > 0 { // thanks to @npawelek https://github.com/npawelek
			continue // already exists
			// TODO remove it instead and continue execution below
		}
		r.logger.Info("adding %s as route via %s", subnetStr, defaultInterface)
		output, err = r.commander.Run("ip", "route", "add", subnetStr, "via", defaultGateway.String(), "dev", defaultInterface)
		if err != nil {
			return fmt.Errorf("cannot add route for %s via %s %s %s: %s: %w", subnetStr, defaultGateway.String(), "dev", defaultInterface, output, err)
		}
	}
	return nil
}
