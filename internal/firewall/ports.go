package firewall

import (
	"context"
	"fmt"
	"strconv"
)

func (c *Config) SetAllowedPort(ctx context.Context, port uint16, intf string) (err error) {
    c.stateMutex.Lock()
    defer c.stateMutex.Unlock()

    interfaceSet, ok := c.allowedInputPorts[port]
    if !ok {
        interfaceSet = make(map[string]struct{})
        c.allowedInputPorts[port] = interfaceSet
    }

    _, alreadySet := interfaceSet[intf]
    if alreadySet {
        return nil
    }

    if c.enabled {
        const remove = false
        err = c.acceptInputToPort(ctx, intf, port, remove)
        if err != nil {
            return fmt.Errorf("accepting input port %d on interface %s: %w",
                port, intf, err)
        }
        
        // ADD THIS: Re-apply user post-rules after port changes
        if err = c.applyUserPostRules(ctx); err != nil {
            return fmt.Errorf("re-applying user post-rules after port change: %w", err)
        }
    }

    interfaceSet[intf] = struct{}{}
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
		err := c.acceptInputToPort(ctx, netInterface, port, remove)
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
