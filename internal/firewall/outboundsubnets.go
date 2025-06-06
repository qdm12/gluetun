package firewall

import (
	"context"
	"fmt"
	"net/netip"

	"github.com/qdm12/gluetun/internal/netlink"
	"github.com/qdm12/gluetun/internal/subnet"
)

func (c *Config) SetOutboundSubnets(ctx context.Context, outboundSubnets []netip.Prefix) (err error) {
    c.stateMutex.Lock()
    defer c.stateMutex.Unlock()

    if !c.enabled {
        c.outboundSubnets = outboundSubnets
        return nil
    }

    // Remove previous outbound subnet rules
    for _, subnet := range c.outboundSubnets {
        subnetIsIPv6 := subnet.Addr().Is6()
        for _, defaultRoute := range c.defaultRoutes {
            defaultRouteIsIPv6 := defaultRoute.Family == netlink.FamilyV6
            ipFamilyMatch := subnetIsIPv6 == defaultRouteIsIPv6
            if !ipFamilyMatch {
                continue
            }

            const remove = true
            err := c.acceptOutputFromIPToSubnet(ctx, defaultRoute.NetInterface,
                defaultRoute.AssignedIP, subnet, remove)
            if err != nil {
                return err
            }
        }
    }

    c.outboundSubnets = outboundSubnets

    // Add new outbound subnet rules
    if err = c.allowOutboundSubnets(ctx); err != nil {
        return err
    }

    // Re-apply user post-rules after subnet changes
    if err = c.applyUserPostRules(ctx); err != nil {
        return fmt.Errorf("re-applying user post-rules after outbound subnet change: %w", err)
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
			err := c.acceptOutputFromIPToSubnet(ctx, defaultRoute.NetInterface,
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
			err := c.acceptOutputFromIPToSubnet(ctx, defaultRoute.NetInterface,
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
