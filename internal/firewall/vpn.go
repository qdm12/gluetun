package firewall

import (
	"context"
	"fmt"

	"github.com/qdm12/gluetun/internal/models"
)

type VPNConnectionSetter interface {
	SetVPNConnection(ctx context.Context,
		connection models.Connection, vpnIntf string) error
}

func (c *Config) SetVPNConnection(ctx context.Context,
	connection models.Connection, vpnIntf string) (err error) {
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
			c.logger.Error("cannot remove outdated VPN connection through firewall: " + err.Error())
		}
	}
	c.vpnConnection = models.Connection{}

	if c.vpnIntf != "" {
		if err = c.acceptOutputThroughInterface(ctx, c.vpnIntf, remove); err != nil {
			c.logger.Error("cannot remove outdated VPN interface from firewall: " + err.Error())
		}
	}
	c.vpnIntf = ""

	remove = false

	if err := c.acceptOutputTrafficToVPN(ctx, c.defaultInterface, connection, remove); err != nil {
		return fmt.Errorf("cannot set VPN connection through firewall: %w", err)
	}
	c.vpnConnection = connection

	if err = c.acceptOutputThroughInterface(ctx, vpnIntf, remove); err != nil {
		return fmt.Errorf("cannot accept output traffic through interface %s: %w", vpnIntf, err)
	}
	c.vpnIntf = vpnIntf

	return nil
}
