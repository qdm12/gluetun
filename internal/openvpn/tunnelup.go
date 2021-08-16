package openvpn

import (
	"context"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/version"
)

func (l *Loop) onTunnelUp(ctx context.Context) {
	vpnDestination, err := l.routing.VPNDestinationIP()
	if err != nil {
		l.logger.Warn(err.Error())
	} else {
		l.logger.Info("VPN routing IP address: " + vpnDestination.String())
	}

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
}
