package openvpn

import (
	"context"
	"time"

	"github.com/qdm12/gluetun/internal/portforward"
	"github.com/qdm12/gluetun/internal/provider"
)

func (l *Loop) startPortForwarding(ctx context.Context,
	portForwarder provider.PortForwarder, serverName string) {
	if !l.GetSettings().Provider.PortForwarding.Enabled {
		return
	}
	// only used for PIA for now
	gateway, err := l.routing.VPNLocalGatewayIP()
	if err != nil {
		l.logger.Error("cannot obtain VPN local gateway IP: " + err.Error())
		return
	}
	l.logger.Info("VPN gateway IP address: " + gateway.String())
	pfData := portforward.StartData{
		PortForwarder: portForwarder,
		Gateway:       gateway,
		ServerName:    serverName,
	}
	_, err = l.portForward.Start(ctx, pfData)
	if err != nil {
		l.logger.Error("cannot start port forwarding: " + err.Error())
	}
}

func (l *Loop) stopPortForwarding(ctx context.Context, timeout time.Duration) {
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	_, err := l.portForward.Stop(ctx)
	if err != nil {
		l.logger.Error("cannot stop port forwarding: " + err.Error())
	}
}
