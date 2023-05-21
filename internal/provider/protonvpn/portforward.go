package protonvpn

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/netip"
	"time"

	"github.com/qdm12/gluetun/internal/natpmp"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

var (
	ErrGatewayIPNotValid = errors.New("gateway IP address is not valid")
)

// PortForward obtains a VPN server side port forwarded from ProtonVPN gateway.
func (p *Provider) PortForward(ctx context.Context, _ *http.Client,
	logger utils.Logger, gateway netip.Addr, _ string) (
	port uint16, err error) {
	if !gateway.IsValid() {
		return 0, fmt.Errorf("%w", ErrGatewayIPNotValid)
	}

	client := natpmp.New()
	_, externalIPv4Address, err := client.ExternalAddress(ctx,
		gateway)
	if err != nil {
		return 0, fmt.Errorf("getting external IPv4 address: %w", err)
	}

	logger.Info("gateway external IPv4 address is " + externalIPv4Address.String())
	const networkProtocol = "udp"
	const internalPort, externalPort = 0, 0
	const lifetime = 60 * time.Second
	_, assignedInternalPort, assignedExternalPort, assignedLiftetime, err :=
		client.AddPortMapping(ctx, gateway, networkProtocol,
			internalPort, externalPort, lifetime)
	if err != nil {
		return 0, fmt.Errorf("adding port mapping: %w", err)
	}

	if assignedLiftetime != lifetime {
		logger.Warn(fmt.Sprintf("assigned lifetime %s differs"+
			" from requested lifetime %s",
			assignedLiftetime, lifetime))
	}

	if assignedInternalPort != assignedExternalPort {
		logger.Warn(fmt.Sprintf("internal port assigned %d differs"+
			" from external port assigned %d",
			assignedInternalPort, assignedExternalPort))
	}

	port = assignedExternalPort
	return port, nil
}

func (p *Provider) KeepPortForward(ctx context.Context, port uint16,
	gateway netip.Addr, _ string, logger utils.Logger) (err error) {
	logger.Info(fmt.Sprintf("keeping port forward with port %d", port))
	client := natpmp.New()
	const refreshTimeout = 45 * time.Second
	timer := time.NewTimer(refreshTimeout)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
		}
		logger.Info("keeping port forward triggered by 45s timer")

		const networkProtocol = "udp"
		const internalPort = 0
		const lifetime = 60 * time.Second
		_, assignedInternalPort, assignedExternalPort, assignedLiftetime, err :=
			client.AddPortMapping(ctx, gateway, networkProtocol,
				internalPort, port, lifetime)
		if err != nil {
			return fmt.Errorf("adding port mapping: %w", err)
		}

		if assignedLiftetime != lifetime {
			logger.Warn(fmt.Sprintf("assigned lifetime %s differs"+
				" from requested lifetime %s",
				assignedLiftetime, lifetime))
		}

		if assignedInternalPort != assignedExternalPort {
			logger.Warn(fmt.Sprintf("internal port assigned %d differs"+
				" from external port assigned %d",
				assignedInternalPort, assignedExternalPort))
		}
		logger.Info(fmt.Sprintf(
			"port %d, assigned internal port %d, assigned external port %d, assignled lifetime %s",
			port, assignedInternalPort, assignedExternalPort, assignedLiftetime))

		timer.Reset(refreshTimeout)
	}
}
