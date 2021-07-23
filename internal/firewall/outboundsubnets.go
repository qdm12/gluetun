package firewall

import (
	"context"
	"fmt"
	"net"
)

func (c *configurator) SetOutboundSubnets(ctx context.Context, subnets []net.IPNet) (err error) {
	c.stateMutex.Lock()
	defer c.stateMutex.Unlock()

	if !c.enabled {
		c.logger.Info("firewall disabled, only updating allowed subnets internal list")
		c.outboundSubnets = make([]net.IPNet, len(subnets))
		copy(c.outboundSubnets, subnets)
		return nil
	}

	c.logger.Info("setting allowed subnets through firewall...")

	subnetsToAdd := findSubnetsToAdd(c.outboundSubnets, subnets)
	subnetsToRemove := findSubnetsToRemove(c.outboundSubnets, subnets)
	if len(subnetsToAdd) == 0 && len(subnetsToRemove) == 0 {
		return nil
	}

	c.removeOutboundSubnets(ctx, subnetsToRemove)
	if err := c.addOutboundSubnets(ctx, subnetsToAdd); err != nil {
		return fmt.Errorf("cannot set allowed subnets through firewall: %w", err)
	}

	return nil
}

func (c *configurator) removeOutboundSubnets(ctx context.Context, subnets []net.IPNet) {
	const remove = true
	for _, subnet := range subnets {
		if err := c.acceptOutputFromIPToSubnet(ctx, c.defaultInterface, c.localIP, subnet, remove); err != nil {
			c.logger.Error("cannot remove outdated outbound subnet through firewall: " + err.Error())
			continue
		}
		c.outboundSubnets = removeSubnetFromSubnets(c.outboundSubnets, subnet)
	}
}

func (c *configurator) addOutboundSubnets(ctx context.Context, subnets []net.IPNet) error {
	const remove = false
	for _, subnet := range subnets {
		if err := c.acceptOutputFromIPToSubnet(ctx, c.defaultInterface, c.localIP, subnet, remove); err != nil {
			return fmt.Errorf("cannot add allowed subnet through firewall: %w", err)
		}
		c.outboundSubnets = append(c.outboundSubnets, subnet)
	}
	return nil
}
