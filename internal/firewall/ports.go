package firewall

import (
	"context"
	"fmt"

	"github.com/qdm12/private-internet-access-docker/internal/constants"
)

func (c *configurator) SetAllowedPort(ctx context.Context, port uint16) (err error) {
	c.stateMutex.Lock()
	defer c.stateMutex.Unlock()

	if port == 0 {
		return nil
	}

	if !c.enabled {
		c.logger.Info("firewall disabled, only updating allowed ports internal list")
		c.allowedPorts[port] = struct{}{}
		return nil
	}

	c.logger.Info("setting allowed port %d through firewall...", port)

	if _, ok := c.allowedPorts[port]; ok {
		return nil
	}

	const remove = false
	if err := c.acceptInputToPort(ctx, "*", constants.TCP, port, remove); err != nil {
		return fmt.Errorf("cannot set allowed port %d through firewall: %w", port, err)
	}
	if err := c.acceptInputToPort(ctx, "*", constants.UDP, port, remove); err != nil {
		return fmt.Errorf("cannot set allowed port %d through firewall: %w", port, err)
	}
	c.allowedPorts[port] = struct{}{}

	return nil
}

func (c *configurator) RemoveAllowedPort(ctx context.Context, port uint16) (err error) {
	c.stateMutex.Lock()
	defer c.stateMutex.Unlock()

	if port == 0 {
		return nil
	}

	if !c.enabled {
		c.logger.Info("firewall disabled, only updating allowed ports internal list")
		delete(c.allowedPorts, port)
		return nil
	}

	c.logger.Info("removing allowed port %d through firewall...", port)

	if _, ok := c.allowedPorts[port]; !ok {
		return nil
	}

	const remove = true
	if err := c.acceptInputToPort(ctx, "*", constants.TCP, port, remove); err != nil {
		return fmt.Errorf("cannot remove allowed port %d through firewall: %w", port, err)
	}
	if err := c.acceptInputToPort(ctx, "*", constants.UDP, port, remove); err != nil {
		return fmt.Errorf("cannot remove allowed port %d through firewall: %w", port, err)
	}
	delete(c.allowedPorts, port)

	return nil
}

// Use 0 to remove
func (c *configurator) SetPortForward(ctx context.Context, port uint16) (err error) {
	c.stateMutex.Lock()
	defer c.stateMutex.Unlock()

	if port == c.portForwarded {
		return nil
	}

	if !c.enabled {
		c.logger.Info("firewall disabled, only updating port forwarded internally")
		c.portForwarded = port
		return nil
	}

	const tun = string(constants.TUN)
	if c.portForwarded > 0 {
		if err := c.acceptInputToPort(ctx, tun, constants.TCP, c.portForwarded, true); err != nil {
			return fmt.Errorf("cannot remove outdated port forward rule from firewall: %w", err)
		}
		if err := c.acceptInputToPort(ctx, tun, constants.UDP, c.portForwarded, true); err != nil {
			return fmt.Errorf("cannot remove outdated port forward rule from firewall: %w", err)
		}
	}

	if port == 0 { // not changing port
		c.portForwarded = 0
		return nil
	}

	if err := c.acceptInputToPort(ctx, tun, constants.TCP, port, false); err != nil {
		return fmt.Errorf("cannot accept port forwarded through firewall: %w", err)
	}
	if err := c.acceptInputToPort(ctx, tun, constants.UDP, port, false); err != nil {
		return fmt.Errorf("cannot accept port forwarded through firewall: %w", err)
	}
	return nil
}
