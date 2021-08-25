package routing

import (
	"errors"
	"fmt"
	"net"

	"github.com/qdm12/gluetun/internal/subnet"
)

var (
	ErrAddOutboundSubnet = errors.New("cannot add outbound subnet to routes")
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
	defaultInterfaceName string, defaultGateway net.IP) error {
	r.stateMutex.Lock()
	defer r.stateMutex.Unlock()

	subnetsToAdd, subnetsToRemove := subnet.FindSubnetsToChange(
		r.outboundSubnets, outboundSubnets)

	if len(subnetsToAdd) == 0 && len(subnetsToRemove) == 0 {
		return nil
	}

	r.removeOutboundSubnets(subnetsToRemove, defaultInterfaceName, defaultGateway)
	return r.addOutboundSubnets(subnetsToAdd, defaultInterfaceName, defaultGateway)
}

func (r *Routing) removeOutboundSubnets(subnets []net.IPNet,
	defaultInterfaceName string, defaultGateway net.IP) {
	for _, subNet := range subnets {
		const table = 0
		if err := r.deleteRouteVia(subNet, defaultGateway, defaultInterfaceName, table); err != nil {
			r.logger.Error("cannot remove outdated outbound subnet from routing: " + err.Error())
			continue
		}
		r.outboundSubnets = subnet.RemoveSubnetFromSubnets(r.outboundSubnets, subNet)
	}
}

func (r *Routing) addOutboundSubnets(subnets []net.IPNet,
	defaultInterfaceName string, defaultGateway net.IP) error {
	for _, subnet := range subnets {
		const table = 0
		if err := r.addRouteVia(subnet, defaultGateway, defaultInterfaceName, table); err != nil {
			return fmt.Errorf("%w: %s: %s", ErrAddOutboundSubnet, subnet, err)
		}
		r.outboundSubnets = append(r.outboundSubnets, subnet)
	}
	return nil
}
