package firewall

import (
	"context"
	"fmt"

	"github.com/qdm12/gluetun/internal/models"
)

func (c *configurator) SetVPNConnections(ctx context.Context, connections []models.OpenVPNConnection) (err error) {
	c.stateMutex.Lock()
	defer c.stateMutex.Unlock()

	if !c.enabled {
		c.logger.Info("firewall disabled, only updating VPN connections internal list")
		c.vpnConnections = make([]models.OpenVPNConnection, len(connections))
		copy(c.vpnConnections, connections)
		return nil
	}

	c.logger.Info("setting VPN connections through firewall...")

	connectionsToAdd := findConnectionsToAdd(c.vpnConnections, connections)
	connectionsToRemove := findConnectionsToRemove(c.vpnConnections, connections)
	if len(connectionsToAdd) == 0 && len(connectionsToRemove) == 0 {
		return nil
	}

	c.removeConnections(ctx, connectionsToRemove, c.defaultInterface)
	if err := c.addConnections(ctx, connectionsToAdd, c.defaultInterface); err != nil {
		return fmt.Errorf("cannot set VPN connections through firewall: %w", err)
	}

	return nil
}

func removeConnectionFromConnections(connections []models.OpenVPNConnection, connection models.OpenVPNConnection) []models.OpenVPNConnection {
	L := len(connections)
	for i := range connections {
		if connection.Equal(connections[i]) {
			connections[i] = connections[L-1]
			connections = connections[:L-1]
			break
		}
	}
	return connections
}

func findConnectionsToAdd(oldConnections, newConnections []models.OpenVPNConnection) (connectionsToAdd []models.OpenVPNConnection) {
	for _, newConnection := range newConnections {
		found := false
		for _, oldConnection := range oldConnections {
			if oldConnection.Equal(newConnection) {
				found = true
				break
			}
		}
		if !found {
			connectionsToAdd = append(connectionsToAdd, newConnection)
		}
	}
	return connectionsToAdd
}

func findConnectionsToRemove(oldConnections, newConnections []models.OpenVPNConnection) (connectionsToRemove []models.OpenVPNConnection) {
	for _, oldConnection := range oldConnections {
		found := false
		for _, newConnection := range newConnections {
			if oldConnection.Equal(newConnection) {
				found = true
				break
			}
		}
		if !found {
			connectionsToRemove = append(connectionsToRemove, oldConnection)
		}
	}
	return connectionsToRemove
}

func (c *configurator) removeConnections(ctx context.Context, connections []models.OpenVPNConnection, defaultInterface string) {
	for _, conn := range connections {
		const remove = true
		if err := c.acceptOutputTrafficToVPN(ctx, defaultInterface, conn, remove); err != nil {
			c.logger.Error("cannot remove outdated VPN connection through firewall: %s", err)
			continue
		}
		c.vpnConnections = removeConnectionFromConnections(c.vpnConnections, conn)
	}
}

func (c *configurator) addConnections(ctx context.Context, connections []models.OpenVPNConnection, defaultInterface string) error {
	const remove = false
	for _, conn := range connections {
		if err := c.acceptOutputTrafficToVPN(ctx, defaultInterface, conn, remove); err != nil {
			return err
		}
		c.vpnConnections = append(c.vpnConnections, conn)
	}
	return nil
}
