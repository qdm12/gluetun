package routing

import (
	"fmt"
	"net"

	"github.com/qdm12/gluetun/internal/subnet"
)

const (
	outboundTable    = 199
	outboundPriority = 99
)

func (r *Routing) SetOutboundRoutes(outboundSubnets []net.IPNet) error {
	defaultRoutes, err := r.DefaultRoutes()
	if err != nil {
		return err
	}
	return r.setOutboundRoutes(outboundSubnets, defaultRoutes)
}

func (r *Routing) setOutboundRoutes(outboundSubnets []net.IPNet,
	defaultRoutes []DefaultRoute) (err error) {
	r.stateMutex.Lock()
	defer r.stateMutex.Unlock()

	subnetsToAdd, subnetsToRemove := subnet.FindSubnetsToChange(
		r.outboundSubnets, outboundSubnets)

	if len(subnetsToAdd) == 0 && len(subnetsToRemove) == 0 {
		return nil
	}

	warnings := r.removeOutboundSubnets(subnetsToRemove, defaultRoutes)
	for _, warning := range warnings {
		r.logger.Warn("cannot remove outdated outbound subnet from routing: " + warning)
	}

	err = r.addOutboundSubnets(subnetsToAdd, defaultRoutes)
	if err != nil {
		return fmt.Errorf("cannot add outbound subnet to routes: %w", err)
	}

	return nil
}

func (r *Routing) removeOutboundSubnets(subnets []net.IPNet,
	defaultRoutes []DefaultRoute) (warnings []string) {
	for i, subNet := range subnets {
		for _, defaultRoute := range defaultRoutes {
			err := r.deleteRouteVia(subNet, defaultRoute.Gateway, defaultRoute.NetInterface, outboundTable)
			if err != nil {
				warnings = append(warnings, err.Error())
				continue
			}
		}

		ruleSrcNet := (*net.IPNet)(nil)
		ruleDstNet := &subnets[i]
		err := r.deleteIPRule(ruleSrcNet, ruleDstNet, outboundTable, outboundPriority)
		if err != nil {
			warnings = append(warnings,
				"cannot delete rule: for subnet "+subNet.String()+": "+err.Error())
			continue
		}

		r.outboundSubnets = subnet.RemoveSubnetFromSubnets(r.outboundSubnets, subNet)
	}

	return warnings
}

func (r *Routing) addOutboundSubnets(subnets []net.IPNet,
	defaultRoutes []DefaultRoute) (err error) {
	for i, subnet := range subnets {
		for _, defaultRoute := range defaultRoutes {
			err = r.addRouteVia(subnet, defaultRoute.Gateway, defaultRoute.NetInterface, outboundTable)
			if err != nil {
				return fmt.Errorf("cannot add route for subnet %s: %w", subnet, err)
			}
		}

		ruleSrcNet := (*net.IPNet)(nil)
		ruleDstNet := &subnets[i]
		err = r.addIPRule(ruleSrcNet, ruleDstNet, outboundTable, outboundPriority)
		if err != nil {
			return fmt.Errorf("cannot add rule: for subnet %s: %w", subnet, err)
		}

		r.outboundSubnets = append(r.outboundSubnets, subnet)
	}
	return nil
}
