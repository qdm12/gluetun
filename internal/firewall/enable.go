package firewall

import (
	"context"
	"fmt"

	"github.com/qdm12/gluetun/internal/constants"
)

func (c *configurator) SetEnabled(ctx context.Context, enabled bool) (err error) {
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
			return err
		}
		c.enabled = false
		c.logger.Info("disabled successfully")
		return nil
	}

	c.logger.Info("enabling...")

	if err := c.enable(ctx); err != nil {
		return err
	}
	c.enabled = true
	c.logger.Info("enabled successfully")

	return nil
}

func (c *configurator) disable(ctx context.Context) (err error) {
	if err = c.clearAllRules(ctx); err != nil {
		return fmt.Errorf("cannot disable firewall: %w", err)
	}
	if err = c.setAllPolicies(ctx, "ACCEPT"); err != nil {
		return fmt.Errorf("cannot disable firewall: %w", err)
	}
	return nil
}

// To use in defered call when enabling the firewall
func (c *configurator) fallbackToDisabled(ctx context.Context) {
	if ctx.Err() != nil {
		return
	}
	if err := c.SetEnabled(ctx, false); err != nil {
		c.logger.Error(err)
	}
}

func (c *configurator) enable(ctx context.Context) (err error) { //nolint:gocognit
	if err = c.setAllPolicies(ctx, "DROP"); err != nil {
		return fmt.Errorf("cannot enable firewall: %w", err)
	}

	const remove = false

	defer func() {
		if err != nil {
			c.fallbackToDisabled(ctx)
		}
	}()

	// Loopback traffic
	if err = c.acceptInputThroughInterface(ctx, "lo", remove); err != nil {
		return fmt.Errorf("cannot enable firewall: %w", err)
	}
	if err = c.acceptOutputThroughInterface(ctx, "lo", remove); err != nil {
		return fmt.Errorf("cannot enable firewall: %w", err)
	}

	if err = c.acceptEstablishedRelatedTraffic(ctx, remove); err != nil {
		return fmt.Errorf("cannot enable firewall: %w", err)
	}
	for _, conn := range c.vpnConnections {
		if err = c.acceptOutputTrafficToVPN(ctx, c.defaultInterface, conn, remove); err != nil {
			return fmt.Errorf("cannot enable firewall: %w", err)
		}
	}
	if err = c.acceptOutputThroughInterface(ctx, string(constants.TUN), remove); err != nil {
		return fmt.Errorf("cannot enable firewall: %w", err)
	}
	if err := c.acceptInputFromSubnetToSubnet(ctx, "*", c.localSubnet, c.localSubnet, remove); err != nil {
		return fmt.Errorf("cannot enable firewall: %w", err)
	}
	if err := c.acceptOutputFromSubnetToSubnet(ctx, "*", c.localSubnet, c.localSubnet, remove); err != nil {
		return fmt.Errorf("cannot enable firewall: %w", err)
	}
	for _, subnet := range c.allowedSubnets {
		if err := c.acceptInputFromSubnetToSubnet(ctx, c.defaultInterface, subnet, c.localSubnet, remove); err != nil {
			return fmt.Errorf("cannot enable firewall: %w", err)
		}
		if err := c.acceptOutputFromSubnetToSubnet(ctx, c.defaultInterface, c.localSubnet, subnet, remove); err != nil {
			return fmt.Errorf("cannot enable firewall: %w", err)
		}
	}
	// Re-ensure all routes exist
	for _, subnet := range c.allowedSubnets {
		if err := c.routing.AddRouteVia(ctx, subnet, c.defaultGateway, c.defaultInterface); err != nil {
			return fmt.Errorf("cannot enable firewall: %w", err)
		}
	}

	for port, intf := range c.allowedInputPorts {
		if err := c.acceptInputToPort(ctx, intf, port, remove); err != nil {
			return fmt.Errorf("cannot enable firewall: %w", err)
		}
	}

	if err := c.runUserPostRules(ctx, "/iptables/post-rules.txt", remove); err != nil {
		return fmt.Errorf("cannot enable firewall: %w", err)
	}

	return nil
}
