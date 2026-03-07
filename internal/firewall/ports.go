package firewall

import (
	"context"
	"fmt"
	"strconv"
)

func (c *Config) SetAllowedPort(ctx context.Context, port uint16, intf string) (err error) {
	c.stateMutex.Lock()
	defer c.stateMutex.Unlock()

	if port == 0 {
		return nil
	}

	if !c.enabled {
		c.logger.Info("firewall disabled, only updating allowed ports internal state")
		existingInterfaces, ok := c.allowedInputPorts[port]
		if !ok {
			existingInterfaces = make(map[string]struct{})
		}
		existingInterfaces[intf] = struct{}{}
		c.allowedInputPorts[port] = existingInterfaces
		return nil
	}

	netInterfaces, has := c.allowedInputPorts[port]
	if !has {
		netInterfaces = make(map[string]struct{})
	} else if _, exists := netInterfaces[intf]; exists {
		return nil
	}

	c.logger.Info("setting allowed input port " + fmt.Sprint(port) + " through interface " + intf + "...")

	const remove = false
	if err := c.impl.AcceptInputToPort(ctx, intf, port, remove); err != nil {
		return fmt.Errorf("allowing input to port %d through interface %s: %w",
			port, intf, err)
	}
	netInterfaces[intf] = struct{}{}
	c.allowedInputPorts[port] = netInterfaces

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

	c.logger.Info("removing allowed port " + strconv.Itoa(int(port)) + "...")

	interfacesSet, ok := c.allowedInputPorts[port]
	if !ok {
		return nil
	}

	const remove = true
	for netInterface := range interfacesSet {
		err := c.impl.AcceptInputToPort(ctx, netInterface, port, remove)
		if err != nil {
			return fmt.Errorf("removing allowed port %d on interface %s: %w",
				port, netInterface, err)
		}
		delete(interfacesSet, netInterface)
	}

	// All interfaces were removed successfully, so remove the port entry.
	delete(c.allowedInputPorts, port)

	return nil
}
