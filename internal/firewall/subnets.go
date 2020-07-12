package firewall

import (
	"context"
	"fmt"
	"net"
)

func (c *configurator) SetAllowedSubnets(ctx context.Context, subnets []net.IPNet) (err error) {
	c.stateMutex.Lock()
	defer c.stateMutex.Unlock()

	if !c.enabled {
		c.logger.Info("firewall disabled, only updating allowed subnets internal list")
		c.allowedSubnets = make([]net.IPNet, len(subnets))
		copy(c.allowedSubnets, subnets)
		return nil
	}

	c.logger.Info("setting allowed subnets through firewall...")

	subnetsToAdd := findSubnetsToAdd(c.allowedSubnets, subnets)
	subnetsToRemove := findSubnetsToRemove(c.allowedSubnets, subnets)
	if len(subnetsToAdd) == 0 && len(subnetsToRemove) == 0 {
		return nil
	}

	defaultInterface, defaultGateway, err := c.routing.DefaultRoute()
	if err != nil {
		return fmt.Errorf("cannot set allowed subnets through firewall: %w", err)
	}

	c.removeSubnets(ctx, subnetsToRemove, defaultInterface)
	if err := c.addSubnets(ctx, subnetsToAdd, defaultInterface, defaultGateway); err != nil {
		return fmt.Errorf("cannot set allowed subnets through firewall: %w", err)
	}

	return nil
}

func findSubnetsToAdd(oldSubnets, newSubnets []net.IPNet) (subnetsToAdd []net.IPNet) {
	for _, newSubnet := range newSubnets {
		found := false
		for _, oldSubnet := range oldSubnets {
			if subnetsAreEqual(oldSubnet, newSubnet) {
				found = true
				break
			}
		}
		if !found {
			subnetsToAdd = append(subnetsToAdd, newSubnet)
		}
	}
	return subnetsToAdd
}

func findSubnetsToRemove(oldSubnets, newSubnets []net.IPNet) (subnetsToRemove []net.IPNet) {
	for _, oldSubnet := range oldSubnets {
		found := false
		for _, newSubnet := range newSubnets {
			if subnetsAreEqual(oldSubnet, newSubnet) {
				found = true
				break
			}
		}
		if !found {
			subnetsToRemove = append(subnetsToRemove, oldSubnet)
		}
	}
	return subnetsToRemove
}

func subnetsAreEqual(a, b net.IPNet) bool {
	return a.IP.Equal(b.IP) && a.Mask.String() == b.Mask.String()
}

func removeSubnetFromSubnets(subnets []net.IPNet, subnet net.IPNet) []net.IPNet {
	L := len(subnets)
	for i := range subnets {
		if subnetsAreEqual(subnet, subnets[i]) {
			subnets[i] = subnets[L-1]
			subnets = subnets[:L-1]
			break
		}
	}
	return subnets
}

func (c *configurator) removeSubnets(ctx context.Context, subnets []net.IPNet, defaultInterface string) {
	const remove = true
	for _, subnet := range subnets {
		failed := false
		if err := c.acceptInputFromToSubnet(ctx, subnet, defaultInterface, remove); err != nil {
			failed = true
			c.logger.Error("cannot remove outdated allowed subnet through firewall: %s", err)
		}
		if err := c.acceptOutputFromToSubnet(ctx, subnet, defaultInterface, remove); err != nil {
			failed = true
			c.logger.Error("cannot remove outdated allowed subnet through firewall: %s", err)
		}
		if err := c.routing.DeleteRouteVia(ctx, subnet); err != nil {
			failed = true
			c.logger.Error("cannot remove outdated allowed subnet route: %s", err)
		}
		if failed {
			continue
		}
		c.allowedSubnets = removeSubnetFromSubnets(c.allowedSubnets, subnet)
	}
}

func (c *configurator) addSubnets(ctx context.Context, subnets []net.IPNet, defaultInterface string, defaultGateway net.IP) error {
	const remove = false
	for _, subnet := range subnets {
		if err := c.acceptInputFromToSubnet(ctx, subnet, defaultInterface, remove); err != nil {
			return fmt.Errorf("cannot add allowed subnet through firewall: %w", err)
		}
		if err := c.acceptOutputFromToSubnet(ctx, subnet, defaultInterface, remove); err != nil {
			return fmt.Errorf("cannot add allowed subnet through firewall: %w", err)
		}
		if err := c.routing.AddRouteVia(ctx, subnet, defaultGateway, defaultInterface); err != nil {
			return fmt.Errorf("cannot add route for allowed subnet: %w", err)
		}
		c.allowedSubnets = append(c.allowedSubnets, subnet)
	}
	return nil
}
