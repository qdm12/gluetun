package firewall

import (
	"context"
	"fmt"

	"github.com/qdm12/gluetun/internal/models"
)

func (c *Config) SetVPNConnection(ctx context.Context,
	connection models.Connection, vpnIntf string,
) (err error) {
	c.stateMutex.Lock()
	defer c.stateMutex.Unlock()

	if !c.enabled {
		c.logger.Info("firewall disabled, only updating internal VPN connection")
		c.vpnConnection = connection
		return nil
	}

	c.logger.Info("allowing VPN connection...")

	if c.vpnConnection.Equal(connection) {
		return nil
	}

	remove := true
	if c.vpnConnection.IP.IsValid() {
		for _, defaultRoute := range c.defaultRoutes {
			if err := c.impl.AcceptOutputTrafficToVPN(ctx, defaultRoute.NetInterface, c.vpnConnection, remove); err != nil {
				c.logger.Error("cannot remove outdated VPN connection rule: " + err.Error())
			}
		}
	}
	c.vpnConnection = models.Connection{}

	if c.vpnIntf != "" {
		if err = c.impl.AcceptOutputThroughInterface(ctx, c.vpnIntf, remove); err != nil {
			c.logger.Error("cannot remove outdated VPN interface rule: " + err.Error())
		}
	}
	c.vpnIntf = ""

	remove = false

	for _, defaultRoute := range c.defaultRoutes {
		if err := c.impl.AcceptOutputTrafficToVPN(ctx, defaultRoute.NetInterface, connection, remove); err != nil {
			return fmt.Errorf("allowing output traffic through VPN connection: %w", err)
		}
	}
	c.vpnConnection = connection

	if err = c.impl.AcceptOutputThroughInterface(ctx, vpnIntf, remove); err != nil {
		return fmt.Errorf("accepting output traffic through interface %s: %w", vpnIntf, err)
	}
	c.vpnIntf = vpnIntf

	return nil
}
