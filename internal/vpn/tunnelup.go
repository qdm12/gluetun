package vpn

import (
	"context"
	"errors"
	"fmt"
	"net/netip"
	"time"

	"github.com/qdm12/dns/v2/pkg/check"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/pmtud"
	"github.com/qdm12/gluetun/internal/version"
	"github.com/qdm12/log"
)

type tunnelUpData struct {
	// Healthcheck
	serverIP netip.Addr
	// vpnType is used for path MTU discovery to find the protocol overhead.
	// It can be "wireguard" or "openvpn".
	vpnType string
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

	mtuLogger := l.logger.New(log.SetComponent("MTU discovery"))
	err := updateToMaxMTU(ctx, data.vpnIntf, data.vpnType,
		l.netLinker, l.routing, mtuLogger)
	if err != nil {
		mtuLogger.Error(err.Error())
	}

	for _, vpnPort := range l.vpnInputPorts {
		err := l.fw.SetAllowedPort(ctx, vpnPort, data.vpnIntf)
		if err != nil {
			l.logger.Error("cannot allow input port through firewall: " + err.Error())
		}
	}

	icmpTarget := l.healthSettings.ICMPTargetIP
	if icmpTarget.IsUnspecified() {
		icmpTarget = data.serverIP
	}
	l.healthChecker.SetConfig(l.healthSettings.TargetAddress, icmpTarget)

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

var errVPNTypeUnknown = errors.New("unknown VPN type")

func updateToMaxMTU(ctx context.Context, vpnInterface string,
	vpnType string, netlinker NetLinker, routing Routing, logger *log.Logger,
) error {
	logger.Info("finding maximum MTU, this can take up to 4 seconds")

	vpnGatewayIP, err := routing.VPNLocalGatewayIP(vpnInterface)
	if err != nil {
		return fmt.Errorf("getting VPN gateway IP address: %w", err)
	}

	link, err := netlinker.LinkByName(vpnInterface)
	if err != nil {
		return fmt.Errorf("getting VPN interface by name: %w", err)
	}

	originalMTU := link.MTU

	// Note: no point testing for an MTU of 1500, it will never work due to the VPN
	// protocol overhead, so start lower than 1500 according to the protocol used.
	const physicalLinkMTU = 1500
	vpnLinkMTU := physicalLinkMTU
	switch vpnType {
	case "wireguard":
		vpnLinkMTU -= 60 // Wireguard overhead
	case "openvpn":
		vpnLinkMTU -= 41 // OpenVPN overhead
	default:
		return fmt.Errorf("%w: %q", errVPNTypeUnknown, vpnType)
	}

	// Setting the VPN link MTU to 1500 might interrupt the connection until
	// the new MTU is set again, but this is necessary to find the highest valid MTU.
	logger.Debugf("VPN interface %s MTU temporarily set to %d", vpnInterface, vpnLinkMTU)

	err = netlinker.LinkSetMTU(link, vpnLinkMTU)
	if err != nil {
		return fmt.Errorf("setting VPN interface %s MTU to %d: %w", vpnInterface, vpnLinkMTU, err)
	}

	const pingTimeout = time.Second
	vpnLinkMTU, err = pmtud.PathMTUDiscover(ctx, vpnGatewayIP, vpnLinkMTU, pingTimeout, logger)
	switch {
	case err == nil:
		logger.Infof("setting VPN interface %s MTU to maximum valid MTU %d", vpnInterface, vpnLinkMTU)
	case errors.Is(err, pmtud.ErrMTUNotFound) || errors.Is(err, pmtud.ErrICMPNotPermitted):
		vpnLinkMTU = int(originalMTU)
		logger.Infof("reverting VPN interface %s MTU to %d (due to: %s)",
			vpnInterface, originalMTU, err)
	default:
		return fmt.Errorf("path MTU discovering: %w", err)
	}

	err = netlinker.LinkSetMTU(link, vpnLinkMTU)
	if err != nil {
		return fmt.Errorf("setting VPN interface %s MTU to %d: %w", vpnInterface, vpnLinkMTU, err)
	}

	return nil
}
