package vpn

import (
	"context"
	"net/netip"

	"github.com/qdm12/dns/v2/pkg/check"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/version"
)

type tunnelUpData struct {
	// Healthcheck
	serverIP netip.Addr
	// Port forwarding
	vpnIntf        string
	serverName     string // used for PIA
	canPortForward bool   // used for PIA
	username       string // used for PIA
	password       string // used for PIA
	portForwarder  PortForwarder
}

func (l *Loop) onTunnelUp(ctx, loopCtx context.Context, data tunnelUpData) {
	l.client.CloseIdleConnections()

	for _, vpnPort := range l.vpnInputPorts {
		err := l.fw.SetAllowedPort(ctx, vpnPort, data.vpnIntf)
		if err != nil {
			l.logger.Error("cannot allow input port through firewall: " + err.Error())
		}
	}

	l.healthChecker.SetICMPTargetIP(data.serverIP)

	healthErrCh, err := l.healthChecker.Start(ctx)
	l.healthServer.SetError(err)
	if err != nil {
		// Note this restart call must be done in a separate goroutine
		// from the VPN loop goroutine.
		l.restartVPN(loopCtx, err)
		return
	}
	defer func() {
		_ = l.healthChecker.Stop()
	}()

	if *l.dnsLooper.GetSettings().DoT.Enabled {
		_, _ = l.dnsLooper.ApplyStatus(ctx, constants.Running)
	} else {
		err := check.WaitForDNS(ctx, check.Settings{})
		if err != nil {
			l.logger.Error("waiting for DNS to be ready: " + err.Error())
		}
	}

	err = l.publicip.RunOnce(ctx)
	if err != nil {
		l.logger.Error("getting public IP address information: " + err.Error())
	}

	if l.versionInfo {
		l.versionInfo = false // only get the version information once
		message, err := version.GetMessage(ctx, l.buildInfo, l.client)
		if err != nil {
			l.logger.Error("cannot get version information: " + err.Error())
		} else {
			l.logger.Info(message)
		}
	}

	err = l.startPortForwarding(data)
	if err != nil {
		l.logger.Error(err.Error())
	}

	select {
	case <-ctx.Done():
	case healthErr := <-healthErrCh:
		l.healthServer.SetError(healthErr)
		// Note this restart call must be done in a separate goroutine
		// from the VPN loop goroutine.
		l.restartVPN(loopCtx, healthErr)
	}
}

func (l *Loop) restartVPN(ctx context.Context, healthErr error) {
	l.logger.Warnf("restarting VPN because it failed to pass the healthcheck: %s", healthErr)
	l.logger.Info("👉 See https://github.com/qdm12/gluetun-wiki/blob/main/faq/healthcheck.md")
	l.logger.Info("DO NOT OPEN AN ISSUE UNLESS YOU HAVE READ AND TRIED EVERY POSSIBLE SOLUTION")
	_, _ = l.ApplyStatus(ctx, constants.Stopped)
	_, _ = l.ApplyStatus(ctx, constants.Running)
}
