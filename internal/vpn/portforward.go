package vpn

import (
	"context"
	"fmt"
	"time"

	"github.com/qdm12/gluetun/internal/portforward"
)

func (l *Loop) startPortForwarding(ctx context.Context, data tunnelUpData) (err error) {
	if !data.portForwarding {
		return nil
	}

	// only used for PIA for now
	gateway, err := l.routing.VPNLocalGatewayIP(data.vpnIntf)
	if err != nil {
		return fmt.Errorf("cannot obtain VPN local gateway IP for interface %s: %w", data.vpnIntf, err)
	}
	l.logger.Info("VPN gateway IP address: " + gateway.String())

	pfData := portforward.StartData{
		PortForwarder: data.portForwarder,
		Gateway:       gateway,
		ServerName:    data.serverName,
		Interface:     data.vpnIntf,
	}
	_, err = l.portForward.Start(ctx, pfData)
	if err != nil {
		return fmt.Errorf("cannot start port forwarding: %w", err)
	}

	return nil
}

func (l *Loop) stopPortForwarding(ctx context.Context,
	timeout time.Duration) (err error) {
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	_, err = l.portForward.Stop(ctx)
	return err
}
