package routing

import (
	"fmt"
	"net"
)

func (r *routing) SetOutboundRoutes(outboundSubnets []net.IPNet) error {
	defaultInterface, defaultGateway, err := r.DefaultRoute()
	if err != nil {
		return fmt.Errorf("cannot set oubtound subnets in routing: %w", err)
	}
	return r.setOutboundRoutes(outboundSubnets, defaultInterface, defaultGateway)
}

func (r *routing) setOutboundRoutes(outboundSubnets []net.IPNet,
	defaultInterfaceName string, defaultGateway net.IP) error {
	r.stateMutex.Lock()
	defer r.stateMutex.Unlock()

	subnetsToRemove := findSubnetsToRemove(r.outboundSubnets, outboundSubnets)
	subnetsToAdd := findSubnetsToAdd(r.outboundSubnets, outboundSubnets)

	if len(subnetsToAdd) == 0 && len(subnetsToRemove) == 0 {
		return nil
	}

	r.removeOutboundSubnets(subnetsToRemove, defaultInterfaceName, defaultGateway)
	if err := r.addOutboundSubnets(subnetsToAdd, defaultInterfaceName, defaultGateway); err != nil {
		return fmt.Errorf("cannot set outbound subnets in routing: %w", err)
	}

	return nil
}

func (r *routing) removeOutboundSubnets(subnets []net.IPNet,
	defaultInterfaceName string, defaultGateway net.IP) {
	for _, subnet := range subnets {
		const table = 0
		if err := r.deleteRouteVia(subnet, defaultGateway, defaultInterfaceName, table); err != nil {
			r.logger.Error("cannot remove outdated outbound subnet from routing: %s", err)
			continue
		}
		r.outboundSubnets = removeSubnetFromSubnets(r.outboundSubnets, subnet)
	}
}

func (r *routing) addOutboundSubnets(subnets []net.IPNet,
	defaultInterfaceName string, defaultGateway net.IP) error {
	for _, subnet := range subnets {
		const table = 0
		if err := r.addRouteVia(subnet, defaultGateway, defaultInterfaceName, table); err != nil {
			return fmt.Errorf("cannot add outbound subnet %s to routing: %w", subnet, err)
		}
		r.outboundSubnets = append(r.outboundSubnets, subnet)
	}
	return nil
}
