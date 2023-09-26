package protonvpn

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/qdm12/gluetun/internal/natpmp"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

// PortForward obtains a VPN server side port forwarded from ProtonVPN gateway.
func (p *Provider) PortForward(ctx context.Context, objects utils.PortForwardObjects) (
	port uint16, err error) {
	client := natpmp.New()
	_, externalIPv4Address, err := client.ExternalAddress(ctx,
		objects.Gateway)
	if err != nil {
		if strings.HasSuffix(err.Error(), "connection refused") {
			err = fmt.Errorf("%w - make sure you have +pmp at the end of your OpenVPN username", err)
		}
		return 0, fmt.Errorf("getting external IPv4 address: %w", err)
	}

	logger := objects.Logger

	logger.Info("gateway external IPv4 address is " + externalIPv4Address.String())
	const internalPort, externalPort = 0, 0
	const lifetime = 60 * time.Second

	_, _, assignedUDPExternalPort, assignedLifetime, err :=
		client.AddPortMapping(ctx, objects.Gateway, "udp",
			internalPort, externalPort, lifetime)
	if err != nil {
		return 0, fmt.Errorf("adding UDP port mapping: %w", err)
	}
	checkLifetime(logger, "UDP", lifetime, assignedLifetime)

	_, _, assignedTCPExternalPort, assignedLifetime, err :=
		client.AddPortMapping(ctx, objects.Gateway, "tcp",
			internalPort, externalPort, lifetime)
	if err != nil {
		return 0, fmt.Errorf("adding TCP port mapping: %w", err)
	}
	checkLifetime(logger, "TCP", lifetime, assignedLifetime)

	checkExternalPorts(logger, assignedUDPExternalPort, assignedTCPExternalPort)
	port = assignedTCPExternalPort

	p.portForwarded = port

	return port, nil
}

func checkLifetime(logger utils.Logger, protocol string,
	requested, actual time.Duration) {
	if requested != actual {
		logger.Warn(fmt.Sprintf("assigned %s port lifetime %s differs"+
			" from requested lifetime %s", strings.ToUpper(protocol),
			actual, requested))
	}
}

func checkExternalPorts(logger utils.Logger, udpPort, tcpPort uint16) {
	if udpPort != tcpPort {
		logger.Warn(fmt.Sprintf("UDP external port %d differs from TCP external port %d",
			udpPort, tcpPort))
	}
}

var ErrExternalPortChanged = errors.New("external port changed")

func (p *Provider) KeepPortForward(ctx context.Context,
	objects utils.PortForwardObjects) (err error) {
	client := natpmp.New()
	const refreshTimeout = 45 * time.Second
	timer := time.NewTimer(refreshTimeout)
	logger := objects.Logger
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
		}

		objects.Logger.Debug("refreshing port forward since 45 seconds have elapsed")
		networkProtocols := []string{"udp", "tcp"}
		const internalPort = 0
		const lifetime = 60 * time.Second

		for _, networkProtocol := range networkProtocols {
			_, _, assignedExternalPort, assignedLiftetime, err :=
				client.AddPortMapping(ctx, objects.Gateway, networkProtocol,
					internalPort, p.portForwarded, lifetime)
			if err != nil {
				return fmt.Errorf("adding port mapping: %w", err)
			}

			if assignedLiftetime != lifetime {
				logger.Warn(fmt.Sprintf("assigned lifetime %s differs"+
					" from requested lifetime %s",
					assignedLiftetime, lifetime))
			}

			if p.portForwarded != assignedExternalPort {
				return fmt.Errorf("%w: %d changed to %d",
					ErrExternalPortChanged, p.portForwarded, assignedExternalPort)
			}
		}

		objects.Logger.Debug(fmt.Sprintf("port forwarded %d maintained", p.portForwarded))

		timer.Reset(refreshTimeout)
	}
}
