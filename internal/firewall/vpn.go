package firewall

import (
	"context"
	"fmt"

	"github.com/qdm12/gluetun/internal/models"
)

func (c *Config) SetVPNConnection(ctx context.Context, connection models.Connection, intf string) (err error) {
    c.stateMutex.Lock()
    defer c.stateMutex.Unlock()

    if !c.enabled {
        c.vpnConnection = connection
        c.vpnIntf = intf
        return nil
    }

    // Remove previous VPN rules
    if c.vpnConnection.IP.IsValid() {
        const remove = true
        interfacesSeen := make(map[string]struct{}, len(c.defaultRoutes))
        for _, defaultRoute := range c.defaultRoutes {
            _, seen := interfacesSeen[defaultRoute.NetInterface]
            if seen {
                continue
            }
            interfacesSeen[defaultRoute.NetInterface] = struct{}{}
            err = c.acceptOutputTrafficToVPN(ctx, defaultRoute.NetInterface, c.vpnConnection, remove)
            if err != nil {
                return fmt.Errorf("removing output traffic through VPN: %w", err)
            }
        }
    }

    c.vpnConnection = connection
    c.vpnIntf = intf

    // Add new VPN rules
    if err = c.allowVPNIP(ctx); err != nil {
        return err
    }

    // Re-apply user post-rules after VPN changes
    if err = c.applyUserPostRules(ctx); err != nil {
        return fmt.Errorf("re-applying user post-rules after VPN change: %w", err)
    }

    return nil
}
