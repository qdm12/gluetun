package firewall

import (
	"context"
	"fmt"

	"github.com/qdm12/gluetun/internal/constants"
)

func (c *configurator) SetVPNPortRedirection(ctx context.Context, src, dst uint16) (err error) {
	c.stateMutex.Lock()
	defer c.stateMutex.Unlock()

	if !c.enabled {
		c.logger.Info("firewall disabled, only updating VPN port redirection internal state")
		c.vpnPortRedirection[src] = dst
		return nil
	}

	c.logger.Info("setting VPN port redirection from port %d to port %d...", src, dst)

	const intf = string(constants.TUN)

	if existingDst, ok := c.vpnPortRedirection[src]; ok {
		if dst == existingDst {
			return nil
		}
		const remove = true
		if err := c.redirectPortToPort(ctx, src, existingDst, intf, remove); err != nil {
			return fmt.Errorf("cannot remove old port redirection from port %d to port %d: %w", src, existingDst, err)
		}
	}

	const remove = false
	if err := c.redirectPortToPort(ctx, src, dst, intf, remove); err != nil {
		return fmt.Errorf("cannot set port redirection from port %d to port %d: %w", src, dst, err)
	}
	c.vpnPortRedirection[src] = dst

	return nil
}

func (c *configurator) RemoveVPNPortRedirection(ctx context.Context, src uint16) (err error) {
	c.stateMutex.Lock()
	defer c.stateMutex.Unlock()

	if !c.enabled {
		c.logger.Info("firewall disabled, only updating VPN port redirection internal list")
		delete(c.vpnPortRedirection, src)
		return nil
	}

	c.logger.Info("removing port redirection from port %d...", src)

	dst, ok := c.vpnPortRedirection[src]
	if !ok {
		return nil
	}

	const remove = true
	const intf = string(constants.TUN)
	if err := c.redirectPortToPort(ctx, src, dst, intf, remove); err != nil {
		return fmt.Errorf("cannot remove port redirection from port %d to port %d: %w", src, dst, err)
	}
	delete(c.vpnPortRedirection, src)

	return nil
}
