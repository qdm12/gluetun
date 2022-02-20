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

type OutboundRoutesSetter interface {
	SetOutboundRoutes(outboundSubnets []net.IPNet) error
}

func (r *Routing) SetOutboundRoutes(outboundSubnets []net.IPNet) error {
	defaultInterface, defaultGateway, err := r.DefaultRoute()
	if err != nil {
		return err
	}
	return r.setOutboundRoutes(outboundSubnets, defaultInterface, defaultGateway)
}

func (r *Routing) setOutboundRoutes(outboundSubnets []net.IPNet,
	defaultInterfaceName string, defaultGateway net.IP) (err error) {
	r.stateMutex.Lock()
	defer r.stateMutex.Unlock()

	subnetsToAdd, subnetsToRemove := subnet.FindSubnetsToChange(
		r.outboundSubnets, outboundSubnets)

	if len(subnetsToAdd) == 0 && len(subnetsToRemove) == 0 {
		return nil
	}

	warnings := r.removeOutboundSubnets(subnetsToRemove, defaultInterfaceName, defaultGateway)
	for _, warning := range warnings {
		r.logger.Warn("cannot remove outdated outbound subnet from routing: " + warning)
	}

	err = r.addOutboundSubnets(subnetsToAdd, defaultInterfaceName, defaultGateway)
	if err != nil {
		return fmt.Errorf("cannot add outbound subnet to routes: %w", err)
	}

	return nil
}

func (r *Routing) removeOutboundSubnets(subnets []net.IPNet,
	defaultInterfaceName string, defaultGateway net.IP) (warnings []string) {
	for i, subNet := range subnets {
		err := r.deleteRouteVia(subNet, defaultGateway, defaultInterfaceName, outboundTable)
		if err != nil {
			warnings = append(warnings, err.Error())
			continue
		}

		ruleSrcNet := (*net.IPNet)(nil)
		ruleDstNet := &subnets[i]
		err = r.deleteIPRule(ruleSrcNet, ruleDstNet, outboundTable, outboundPriority)
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
	defaultInterfaceName string, defaultGateway net.IP) error {
	for i, subnet := range subnets {
		err := r.addRouteVia(subnet, defaultGateway, defaultInterfaceName, outboundTable)
		if err != nil {
			return fmt.Errorf("cannot add route for subnet %s: %w", subnet, err)
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
