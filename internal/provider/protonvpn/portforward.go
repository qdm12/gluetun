package protonvpn

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"strings"
	"time"

	"github.com/qdm12/gluetun/internal/natpmp"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

var ErrServerPortForwardNotSupported = errors.New("server does not support port forwarding")

const nonSymmetricPortStart uint16 = 56789

// PortForward obtains a VPN server side port forwarded from ProtonVPN gateway.
func (p *Provider) PortForward(ctx context.Context, objects utils.PortForwardObjects) (
	internalToExternalPorts map[uint16]uint16, err error,
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

	logger.Debug("gateway external IPv4 address is " + externalIPv4Address.String())
	const externalPort = 0
	const lifetime = 60 * time.Second

	p.internalToExternalPorts = make(map[uint16]uint16, objects.PortsCount)
	for i := range objects.PortsCount {
		internalPort := nonSymmetricPortStart + i
		protoToInternalPort := map[string]uint16{
			"udp": 0,
			"tcp": 0,
		}
		protoToExternalPort := maps.Clone(protoToInternalPort)
		for protocol := range protoToExternalPort {
			_, assignedInternalPort, assignedExternalPort, assignedLifetime, err := client.AddPortMapping(
				ctx, objects.Gateway, protocol, internalPort, externalPort, lifetime)
			if err != nil {
				return nil, fmt.Errorf("adding %d/%d %s port mapping: %w",
					i+1, objects.PortsCount, strings.ToUpper(protocol), err)
			}
			checkLifetime(logger, strings.ToUpper(protocol), lifetime, assignedLifetime)
			checkInternalPort(logger, internalPort, assignedInternalPort)
			protoToInternalPort[protocol] = assignedInternalPort
			protoToExternalPort[protocol] = assignedExternalPort
		}

		checkInternalPorts(logger, protoToInternalPort["udp"], protoToInternalPort["tcp"])
		checkExternalPorts(logger, protoToExternalPort["udp"], protoToExternalPort["tcp"])
		p.internalToExternalPorts[protoToInternalPort["tcp"]] = protoToExternalPort["tcp"]
	}

	return maps.Clone(p.internalToExternalPorts), nil
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

func checkInternalPort(logger utils.Logger, sent, received uint16) {
	if sent != received {
		logger.Warn(fmt.Sprintf("internal port assigned %d differs from requested internal port %d",
			sent, received))
	}
}

func checkInternalPorts(logger utils.Logger, udpPort, tcpPort uint16) {
	if udpPort != tcpPort {
		logger.Warn(fmt.Sprintf("UDP internal port %d differs from TCP internal port %d",
			udpPort, tcpPort))
	}
}

func checkExternalPorts(logger utils.Logger, udpPort, tcpPort uint16) {
	if udpPort != tcpPort {
		logger.Warn(fmt.Sprintf("UDP external port %d differs from TCP external port %d",
			udpPort, tcpPort))
	}
}

var (
	ErrInternalPortChanged = errors.New("internal port changed")
	ErrExternalPortChanged = errors.New("external port changed")
)

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
		networkProtocols := [...]string{"udp", "tcp"}
		const lifetime = 60 * time.Second
		for internalPort, externalPort := range p.internalToExternalPorts {
			for _, networkProtocol := range networkProtocols {
				_, assignedInternalPort, assignedExternalPort, assignedLiftetime, err := client.AddPortMapping(
					ctx, objects.Gateway, networkProtocol, internalPort, externalPort, lifetime)
				if err != nil {
					return fmt.Errorf("adding port mapping: %w", err)
				}
				checkLifetime(logger, networkProtocol, lifetime, assignedLiftetime)
				if externalPort != assignedExternalPort {
					return fmt.Errorf("%w: %d changed to %d",
						ErrExternalPortChanged, externalPort, assignedExternalPort)
				} else if internalPort != assignedInternalPort {
					return fmt.Errorf("%w: %d (for external port %d) changed to %d",
						ErrInternalPortChanged, internalPort, externalPort, assignedInternalPort)
				}
			}
			objects.Logger.Debug(fmt.Sprintf("port forwarded %d maintained", externalPort))
		}

		timer.Reset(refreshTimeout)
	}
}
