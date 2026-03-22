package protonvpn

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/qdm12/gluetun/internal/natpmp"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

var ErrServerPortForwardNotSupported = errors.New("server does not support port forwarding")

// PortForward obtains a VPN server side port forwarded from ProtonVPN gateway.
func (p *Provider) PortForward(ctx context.Context, objects utils.PortForwardObjects) (
	ports []uint16, err error,
) {
	if !objects.CanPortForward {
		return nil, fmt.Errorf("%w", ErrServerPortForwardNotSupported)
	}

	client := natpmp.New()
	_, externalIPv4Address, err := client.ExternalAddress(ctx, objects.Gateway)
	if err != nil {
		switch {
		case strings.HasSuffix(err.Error(), "connection refused"):
			err = fmt.Errorf("%w - make sure you have +pmp at the end of your OpenVPN username "+
				"or that your Wireguard key is set to work with PMP", err)
		case strings.Contains(err.Error(), "i/o timeout"):
			err = fmt.Errorf("%w - make sure FIREWALL_OUTBOUND_SUBNETS does not conflict with "+
				"the VPN gateway ip address %s", err, objects.Gateway)
		}
		return nil, fmt.Errorf("getting external IPv4 address: %w", err)
	}

	logger := objects.Logger

	logger.Info("gateway external IPv4 address is " + externalIPv4Address.String())
	const externalPort = 1
	const lifetime = 60 * time.Second

	p.portsForwarded = make([]uint16, objects.PortsCount)
	for i := range p.portsForwarded {
		internalPort := uint16(i + 1) //nolint:gosec
		protoToPort := map[string]uint16{
			"udp": 0,
			"tcp": 0,
		}
		for protocol := range protoToPort {
			_, _, assignedExternalPort, assignedLifetime, err := client.AddPortMapping(ctx, objects.Gateway, protocol,
				internalPort, externalPort, lifetime)
			if err != nil {
				return nil, fmt.Errorf("adding %d/%d %s port mapping: %w",
					i+1, len(p.portsForwarded), strings.ToUpper(protocol), err)
			}
			checkLifetime(logger, strings.ToUpper(protocol), lifetime, assignedLifetime)
			protoToPort[protocol] = assignedExternalPort
		}

		checkExternalPorts(logger, protoToPort["udp"], protoToPort["tcp"])
		p.portsForwarded[i] = protoToPort["tcp"] // use TCP port as the forwarded port, UDP is the same as TCP
	}

	return slices.Clone(p.portsForwarded), nil
}

func checkLifetime(logger utils.Logger, protocol string,
	requested, actual time.Duration,
) {
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
	objects utils.PortForwardObjects,
) (err error) {
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

		objects.Logger.Debug("refreshing forwarded ports since 45 seconds have elapsed")
		networkProtocols := []string{"udp", "tcp"}
		const lifetime = 60 * time.Second

		for i, portForwarded := range p.portsForwarded {
			internalPort := uint16(i + 1) //nolint:gosec
			for _, networkProtocol := range networkProtocols {
				_, _, assignedExternalPort, assignedLiftetime, err := client.AddPortMapping(ctx, objects.Gateway, networkProtocol,
					internalPort, portForwarded, lifetime)
				if err != nil {
					return fmt.Errorf("adding port mapping: %w", err)
				}

				if assignedLiftetime != lifetime {
					logger.Warn(fmt.Sprintf("assigned lifetime %s differs"+
						" from requested lifetime %s", assignedLiftetime, lifetime))
				}

				if portForwarded != assignedExternalPort {
					return fmt.Errorf("%w: %d changed to %d",
						ErrExternalPortChanged, portForwarded, assignedExternalPort)
				}
				objects.Logger.Debug(fmt.Sprintf("port forwarded %d maintained", portForwarded))
			}
		}

		timer.Reset(refreshTimeout)
	}
}
