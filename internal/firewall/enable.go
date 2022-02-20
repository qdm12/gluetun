package firewall

import (
	"context"
	"fmt"
)

type Enabler interface {
	SetEnabled(ctx context.Context, enabled bool) (err error)
}

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
			return fmt.Errorf("cannot disable firewall: %w", err)
		}
		c.enabled = false
		c.logger.Info("disabled successfully")
		return nil
	}

	c.logger.Info("enabling...")

	if err := c.enable(ctx); err != nil {
		return fmt.Errorf("cannot enable firewall: %w", err)
	}
	c.enabled = true
	c.logger.Info("enabled successfully")

	return nil
}

func (c *Config) disable(ctx context.Context) (err error) {
	if err = c.clearAllRules(ctx); err != nil {
		return fmt.Errorf("cannot clear all rules: %w", err)
	}
	if err = c.setIPv4AllPolicies(ctx, "ACCEPT"); err != nil {
		return fmt.Errorf("cannot set ipv4 policies: %w", err)
	}
	if err = c.setIPv6AllPolicies(ctx, "ACCEPT"); err != nil {
		return fmt.Errorf("cannot set ipv6 policies: %w", err)
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
	if c.vpnConnection.IP != nil {
		if err = c.acceptOutputTrafficToVPN(ctx, c.defaultInterface, c.vpnConnection, remove); err != nil {
			return err
		}
		if err = c.acceptOutputThroughInterface(ctx, c.vpnIntf, remove); err != nil {
			return err
		}
	}

	for _, network := range c.localNetworks {
		if err := c.acceptOutputFromIPToSubnet(ctx, network.InterfaceName, network.IP, *network.IPNet, remove); err != nil {
			return err
		}
	}

	for _, subnet := range c.outboundSubnets {
		if err := c.acceptOutputFromIPToSubnet(ctx, c.defaultInterface, c.localIP, subnet, remove); err != nil {
			return err
		}
	}

	// Allows packets from any IP address to go through eth0 / local network
	// to reach Gluetun.
	for _, network := range c.localNetworks {
		if err := c.acceptInputToSubnet(ctx, network.InterfaceName, *network.IPNet, remove); err != nil {
			return err
		}
	}

	for port, intf := range c.allowedInputPorts {
		if err := c.acceptInputToPort(ctx, intf, port, remove); err != nil {
			return err
		}
	}

	if err := c.runUserPostRules(ctx, c.customRulesPath, remove); err != nil {
		return fmt.Errorf("cannot run user defined post firewall rules: %w", err)
	}

	return nil
}
