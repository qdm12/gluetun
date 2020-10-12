package firewall

import (
	"context"
	"fmt"

	"github.com/qdm12/gluetun/internal/models"
)

func (c *configurator) SetVPNConnection(ctx context.Context, connection models.OpenVPNConnection) (err error) {
	c.stateMutex.Lock()
	defer c.stateMutex.Unlock()

	if !c.enabled {
		c.logger.Info("firewall disabled, only updating internal VPN connection")
		c.vpnConnection = connection
		return nil
	}

	c.logger.Info("setting VPN connection through firewall...")

	if c.vpnConnection.Equal(connection) {
		return nil
	}

	remove := true
	if c.vpnConnection.IP != nil {
		if err := c.acceptOutputTrafficToVPN(ctx, c.defaultInterface, c.vpnConnection, remove); err != nil {
			c.logger.Error("cannot remove outdated VPN connection through firewall: %s", err)
		}
	}
	c.vpnConnection = models.OpenVPNConnection{}
	remove = false
	if err := c.acceptOutputTrafficToVPN(ctx, c.defaultInterface, connection, remove); err != nil {
		return fmt.Errorf("cannot set VPN connection through firewall: %w", err)
	}
	c.vpnConnection = connection
	return nil
}
