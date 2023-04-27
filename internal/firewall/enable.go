package firewall

import (
	"context"
	"fmt"
)

func (c *Config) SetEnabled(ctx context.Context, enabled bool) (err error) {
	c.stateMutex.Lock()
	defer c.stateMutex.Unlock()

	if enabled == c.enabled {
		if enabled {
			c.logger.Info("already enabled")
		} else {
			c.logger.Info("already disabled")
		}
		return nil
	}

	if !enabled {
		c.logger.Info("disabling...")
		if err = c.disable(ctx); err != nil {
			return fmt.Errorf("disabling firewall: %w", err)
		}
		c.enabled = false
		c.logger.Info("disabled successfully")
		return nil
	}

	c.logger.Info("enabling...")

	if err := c.enable(ctx); err != nil {
		return fmt.Errorf("enabling firewall: %w", err)
	}
	c.enabled = true
	c.logger.Info("enabled successfully")

	return nil
}

func (c *Config) disable(ctx context.Context) (err error) {
	if err = c.clearAllRules(ctx); err != nil {
		return fmt.Errorf("clearing all rules: %w", err)
	}
	if err = c.setIPv4AllPolicies(ctx, "ACCEPT"); err != nil {
		return fmt.Errorf("setting ipv4 policies: %w", err)
	}
	if err = c.setIPv6AllPolicies(ctx, "ACCEPT"); err != nil {
		return fmt.Errorf("setting ipv6 policies: %w", err)
	}
	return nil
}

// To use in defered call when enabling the firewall.
func (c *Config) fallbackToDisabled(ctx context.Context) {
	if ctx.Err() != nil {
		return
	}
	if err := c.disable(ctx); err != nil {
		c.logger.Error("failed reversing firewall changes: " + err.Error())
	}
}

func (c *Config) enable(ctx context.Context) (err error) {
	touched := false
	if err = c.setIPv4AllPolicies(ctx, "DROP"); err != nil {
		return err
	}
	touched = true

	if err = c.setIPv6AllPolicies(ctx, "DROP"); err != nil {
		return err
	}

	const remove = false

	defer func() {
		if touched && err != nil {
			c.fallbackToDisabled(ctx)
		}
	}()

	// Loopback traffic
	if err = c.acceptInputThroughInterface(ctx, "lo", remove); err != nil {
		return err
	}
	if err = c.acceptOutputThroughInterface(ctx, "lo", remove); err != nil {
		return err
	}

	if err = c.acceptEstablishedRelatedTraffic(ctx, remove); err != nil {
		return err
	}

	if err = c.allowVPNIP(ctx); err != nil {
		return err
	}

	for _, network := range c.localNetworks {
		if err := c.acceptOutputFromIPToSubnet(ctx, network.InterfaceName, network.IP, network.IPNet, remove); err != nil {
			return err
		}
		if err = c.acceptIpv6MulticastOutput(ctx, network.InterfaceName, remove); err != nil {
			return err
		}
	}

	if err = c.allowOutboundSubnets(ctx); err != nil {
		return err
	}

	// Allows packets from any IP address to go through eth0 / local network
	// to reach Gluetun.
	for _, network := range c.localNetworks {
		if err := c.acceptInputToSubnet(ctx, network.InterfaceName, network.IPNet, remove); err != nil {
			return err
		}
	}

	if err = c.allowInputPorts(ctx); err != nil {
		return err
	}

	if err := c.runUserPostRules(ctx, c.customRulesPath, remove); err != nil {
		return fmt.Errorf("running user defined post firewall rules: %w", err)
	}

	return nil
}

func (c *Config) allowVPNIP(ctx context.Context) (err error) {
	if c.vpnConnection.IP == nil {
		return nil
	}

	const remove = false
	for _, defaultRoute := range c.defaultRoutes {
		err = c.acceptOutputTrafficToVPN(ctx, defaultRoute.NetInterface, c.vpnConnection, remove)
		if err != nil {
			return fmt.Errorf("accepting output traffic through VPN: %w", err)
		}
	}

	return nil
}

func (c *Config) allowOutboundSubnets(ctx context.Context) (err error) {
	for _, subnet := range c.outboundSubnets {
		for _, defaultRoute := range c.defaultRoutes {
			const remove = false
			err := c.acceptOutputFromIPToSubnet(ctx, defaultRoute.NetInterface,
				defaultRoute.AssignedIP, subnet, remove)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Config) allowInputPorts(ctx context.Context) (err error) {
	for port, netInterfaces := range c.allowedInputPorts {
		for netInterface := range netInterfaces {
			const remove = false
			err = c.acceptInputToPort(ctx, netInterface, port, remove)
			if err != nil {
				return fmt.Errorf("accepting input port %d on interface %s: %w",
					port, netInterface, err)
			}
		}
	}
	return nil
}
