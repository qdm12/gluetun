package firewall

import (
	"context"
	"fmt"
)

func (c *Config) SetAllowedPort(ctx context.Context, port uint16, intf string) (err error) {
	c.stateMutex.Lock()
	defer c.stateMutex.Unlock()

	if !c.enabled {
		c.logger.Info("firewall disabled, only updating allowed ports internal list")
		c.setAllowedPortInternalList(port, intf)
		return nil
	}

	c.logger.Info(fmt.Sprintf("setting allowed port %d...", port))

	interfaces, portExists := c.allowedInputPorts[port]
	if !portExists {
		interfaces = make(map[string]struct{})
		c.allowedInputPorts[port] = interfaces
	}

	_, interfaceExists := interfaces[intf]
	if interfaceExists {
		return nil
	}

	const remove = false
	err = c.acceptInputToPort(ctx, intf, port, remove)
	if err != nil {
		return fmt.Errorf("accepting input to port %d: %w", port, err)
	}

	interfaces[intf] = struct{}{}

	// Apply post-rules once after adding the port
	if err := c.applyPostRulesOnce(ctx); err != nil {
		return fmt.Errorf("applying post firewall rules: %w", err)
	}

	return nil
}

func (c *Config) RemoveAllowedPort(ctx context.Context, port uint16, intf string) (err error) {
	c.stateMutex.Lock()
	defer c.stateMutex.Unlock()

	if !c.enabled {
		c.logger.Info("firewall disabled, only updating allowed ports internal list")
		c.removeAllowedPortInternalList(port, intf)
		return nil
	}

	c.logger.Info(fmt.Sprintf("removing allowed port %d...", port))

	interfaces, portExists := c.allowedInputPorts[port]
	if !portExists {
		return nil
	}

	_, interfaceExists := interfaces[intf]
	if !interfaceExists {
		return nil
	}

	const remove = true
	err = c.acceptInputToPort(ctx, intf, port, remove)
	if err != nil {
		return fmt.Errorf("removing input to port %d: %w", port, err)
	}

	delete(interfaces, intf)
	if len(interfaces) == 0 {
		delete(c.allowedInputPorts, port)
	}

	return nil
}

func (c *Config) setAllowedPortInternalList(port uint16, intf string) {
	interfaces, portExists := c.allowedInputPorts[port]
	if !portExists {
		interfaces = make(map[string]struct{})
		c.allowedInputPorts[port] = interfaces
	}
	interfaces[intf] = struct{}{}
}

func (c *Config) removeAllowedPortInternalList(port uint16, intf string) {
	interfaces, portExists := c.allowedInputPorts[port]
	if !portExists {
		return
	}

	delete(interfaces, intf)
	if len(interfaces) == 0 {
		delete(c.allowedInputPorts, port)
	}
}
