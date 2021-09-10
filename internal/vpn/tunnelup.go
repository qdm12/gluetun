package vpn

import (
	"context"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/provider"
	"github.com/qdm12/gluetun/internal/version"
)

type tunnelUpData struct {
	// Port forwarding
	portForwarding bool
	vpnIntf        string
	serverName     string
	portForwarder  provider.PortForwarder
}

func (l *Loop) onTunnelUp(ctx context.Context, data tunnelUpData) {
	l.client.CloseIdleConnections()

	if l.dnsLooper.GetSettings().Enabled {
		_, _ = l.dnsLooper.ApplyStatus(ctx, constants.Running)
	}

	// Runs the Public IP getter job once
	_, _ = l.publicip.ApplyStatus(ctx, constants.Running)
	if l.versionInfo {
		l.versionInfo = false // only get the version information once
		message, err := version.GetMessage(ctx, l.buildInfo, l.client)
		if err != nil {
			l.logger.Error("cannot get version information: " + err.Error())
		} else {
			l.logger.Info(message)
		}
	}

	err := l.startPortForwarding(ctx, data)
	if err != nil {
		l.logger.Error(err.Error())
	}
}
