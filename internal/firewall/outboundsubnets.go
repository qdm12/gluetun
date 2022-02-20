package firewall

import (
	"context"
	"fmt"
	"net"

	"github.com/qdm12/gluetun/internal/subnet"
)

type OutboundSubnetsSetter interface {
	SetOutboundSubnets(ctx context.Context, subnets []net.IPNet) (err error)
}

func (c *Config) SetOutboundSubnets(ctx context.Context, subnets []net.IPNet) (err error) {
	c.stateMutex.Lock()
	defer c.stateMutex.Unlock()

	if !c.enabled {
		c.logger.Info("firewall disabled, only updating allowed subnets internal list")
		c.outboundSubnets = make([]net.IPNet, len(subnets))
		copy(c.outboundSubnets, subnets)
		return nil
	}

	c.logger.Info("setting allowed subnets...")

	subnetsToAdd, subnetsToRemove := subnet.FindSubnetsToChange(c.outboundSubnets, subnets)
	if len(subnetsToAdd) == 0 && len(subnetsToRemove) == 0 {
		return nil
	}

	c.removeOutboundSubnets(ctx, subnetsToRemove)
	if err := c.addOutboundSubnets(ctx, subnetsToAdd); err != nil {
		return fmt.Errorf("cannot set allowed outbound subnets: %w", err)
	}

	return nil
}

func (c *Config) removeOutboundSubnets(ctx context.Context, subnets []net.IPNet) {
	const remove = true
	for _, subNet := range subnets {
		if err := c.acceptOutputFromIPToSubnet(ctx, c.defaultInterface, c.localIP, subNet, remove); err != nil {
			c.logger.Error("cannot remove outdated outbound subnet: " + err.Error())
			continue
		}
		c.outboundSubnets = subnet.RemoveSubnetFromSubnets(c.outboundSubnets, subNet)
	}
}

func (c *Config) addOutboundSubnets(ctx context.Context, subnets []net.IPNet) error {
	const remove = false
	for _, subnet := range subnets {
		if err := c.acceptOutputFromIPToSubnet(ctx, c.defaultInterface, c.localIP, subnet, remove); err != nil {
			return err
		}
		c.outboundSubnets = append(c.outboundSubnets, subnet)
	}
	return nil
}
