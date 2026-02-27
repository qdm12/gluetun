package firewall

import (
	"context"
	"fmt"
	"net/netip"

	"github.com/qdm12/gluetun/internal/netlink"
	"github.com/qdm12/gluetun/internal/subnet"
)

func (c *Config) SetOutboundSubnets(ctx context.Context, subnets []netip.Prefix) (err error) {
	c.stateMutex.Lock()
	defer c.stateMutex.Unlock()

	if !c.enabled {
		c.logger.Info("firewall disabled, only updating allowed subnets internal list")
		c.outboundSubnets = make([]netip.Prefix, len(subnets))
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
		return fmt.Errorf("setting allowed outbound subnets: %w", err)
	}

	return nil
}

func (c *Config) removeOutboundSubnets(ctx context.Context, subnets []netip.Prefix) {
	const remove = true
	for _, subNet := range subnets {
		subnetIsIPv6 := subNet.Addr().Is6()
		firewallUpdated := false
		for _, defaultRoute := range c.defaultRoutes {
			defaultRouteIsIPv6 := defaultRoute.Family == netlink.FamilyV6
			ipFamilyMatch := subnetIsIPv6 == defaultRouteIsIPv6
			if !ipFamilyMatch {
				continue
			}

			firewallUpdated = true
			err := c.impl.AcceptOutputFromIPToSubnet(ctx, defaultRoute.NetInterface,
				defaultRoute.AssignedIP, subNet, remove)
			if err != nil {
				c.logger.Error("cannot remove outdated outbound subnet: " + err.Error())
				continue
			}
		}

		if !firewallUpdated {
			c.logIgnoredSubnetFamily(subNet)
			continue
		}
		c.outboundSubnets = subnet.RemoveSubnetFromSubnets(c.outboundSubnets, subNet)
	}
}

func (c *Config) addOutboundSubnets(ctx context.Context, subnets []netip.Prefix) error {
	const remove = false
	for _, subnet := range subnets {
		subnetIsIPv6 := subnet.Addr().Is6()
		firewallUpdated := false
		for _, defaultRoute := range c.defaultRoutes {
			defaultRouteIsIPv6 := defaultRoute.Family == netlink.FamilyV6
			ipFamilyMatch := subnetIsIPv6 == defaultRouteIsIPv6
			if !ipFamilyMatch {
				continue
			}

			firewallUpdated = true
			err := c.impl.AcceptOutputFromIPToSubnet(ctx, defaultRoute.NetInterface,
				defaultRoute.AssignedIP, subnet, remove)
			if err != nil {
				return err
			}
		}

		if !firewallUpdated {
			c.logIgnoredSubnetFamily(subnet)
			continue
		}
		c.outboundSubnets = append(c.outboundSubnets, subnet)
	}
	return nil
}
