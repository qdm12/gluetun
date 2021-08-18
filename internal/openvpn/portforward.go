package openvpn

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/portforward"
	"github.com/qdm12/gluetun/internal/provider"
)

var (
	errObtainVPNLocalGateway = errors.New("cannot obtain VPN local gateway IP")
	errStartPortForwarding   = errors.New("cannot start port forwarding")
)

func (l *Loop) startPortForwarding(ctx context.Context, enabled bool,
	portForwarder provider.PortForwarder, serverName string) (err error) {
	if !enabled {
		return nil
	}

	// only used for PIA for now
	gateway, err := l.routing.VPNLocalGatewayIP()
	if err != nil {
		return fmt.Errorf("%w: %s", errObtainVPNLocalGateway, err)
	}
	l.logger.Info("VPN gateway IP address: " + gateway.String())

	pfData := portforward.StartData{
		PortForwarder: portForwarder,
		Gateway:       gateway,
		ServerName:    serverName,
		Interface:     constants.TUN,
	}
	_, err = l.portForward.Start(ctx, pfData)
	if err != nil {
		return fmt.Errorf("%w: %s", errStartPortForwarding, err)
	}

	return nil
}

func (l *Loop) stopPortForwarding(ctx context.Context, enabled bool,
	timeout time.Duration) {
	if !enabled {
		return // nothing to stop
	}

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
