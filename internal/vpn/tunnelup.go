package vpn

import (
	"context"
	"strings"
	"time"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/version"
)

type tunnelUpData struct {
	// Port forwarding
	vpnIntf        string
	serverName     string // used for PIA
	canPortForward bool   // used for PIA
	username       string // used for PIA
	password       string // used for PIA
	portForwarder  PortForwarder
}

func (l *Loop) onTunnelUp(ctx context.Context, data tunnelUpData) {
	l.client.CloseIdleConnections()

	for _, vpnPort := range l.vpnInputPorts {
		err := l.fw.SetAllowedPort(ctx, vpnPort, data.vpnIntf)
		if err != nil {
			l.logger.Error("cannot allow input port through firewall: " + err.Error())
		}
	}

	if *l.dnsLooper.GetSettings().DoT.Enabled {
		_, _ = l.dnsLooper.ApplyStatus(ctx, constants.Running)
	}

	l.fetchPublicIPWithRetry(ctx)

	if l.versionInfo {
		l.versionInfo = false // only get the version information once
		message, err := version.GetMessage(ctx, l.buildInfo, l.client)
		if err != nil {
			l.logger.Error("cannot get version information: " + err.Error())
		} else {
			l.logger.Info(message)
		}
	}

	err := l.startPortForwarding(data)
	if err != nil {
		l.logger.Error(err.Error())
	}
}

func (l *Loop) fetchPublicIPWithRetry(ctx context.Context) {
	for {
		err := l.publicip.RunOnce(ctx)
		if err == nil {
			return
		}

		l.logger.Error("getting public IP address information: " + err.Error())
		if !strings.HasSuffix(err.Error(), "read: connection refused") {
			return
		}

		// Retry mechanism asked in https://github.com/qdm12/gluetun/issues/2325
		const retryPeriod = 2 * time.Second
		l.logger.Infof("retrying public IP address information fetch in %s", retryPeriod)
		timer := time.NewTimer(retryPeriod)
		select {
		case <-ctx.Done():
		case <-timer.C:
		}
	}
}
