package firewall

import (
	"context"
	"fmt"
	"strconv"
)

type PortAllower interface {
	SetAllowedPort(ctx context.Context, port uint16, intf string) (err error)
	RemoveAllowedPort(ctx context.Context, port uint16) (err error)
}

func (c *Config) SetAllowedPort(ctx context.Context, port uint16, intf string) (err error) {
	c.stateMutex.Lock()
	defer c.stateMutex.Unlock()

	if port == 0 {
		return nil
	}

	if !c.enabled {
		c.logger.Info("firewall disabled, only updating allowed ports internal state")
		c.allowedInputPorts[port] = intf
		return nil
	}

	c.logger.Info("setting allowed input port " + fmt.Sprint(port) + " through interface " + intf + "...")

	if existingIntf, ok := c.allowedInputPorts[port]; ok {
		if intf == existingIntf {
			return nil
		}
		const remove = true
		if err := c.acceptInputToPort(ctx, existingIntf, port, remove); err != nil {
			return fmt.Errorf("cannot remove old allowed port %d through interface %s: %w", port, existingIntf, err)
		}
	}

	const remove = false
	if err := c.acceptInputToPort(ctx, intf, port, remove); err != nil {
		return fmt.Errorf("cannot set allowed port %d through interface %s: %w", port, intf, err)
	}
	c.allowedInputPorts[port] = intf

	return nil
}

func (c *Config) RemoveAllowedPort(ctx context.Context, port uint16) (err error) {
	c.stateMutex.Lock()
	defer c.stateMutex.Unlock()

	if port == 0 {
		return nil
	}

	if !c.enabled {
		c.logger.Info("firewall disabled, only updating allowed ports internal list")
		delete(c.allowedInputPorts, port)
		return nil
	}

	c.logger.Info("removing allowed port "+strconv.Itoa(int(port))+" through firewall...", port)

	intf, ok := c.allowedInputPorts[port]
	if !ok {
		return nil
	}

	const remove = true
	if err := c.acceptInputToPort(ctx, intf, port, remove); err != nil {
		return fmt.Errorf("cannot remove allowed port %d through interface %s: %w", port, intf, err)
	}
	delete(c.allowedInputPorts, port)

	return nil
}
