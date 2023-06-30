package utils

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/netip"
)

type NoPortForwarder interface {
	PortForward(ctx context.Context, client *http.Client,
		logger Logger, gateway netip.Addr, serverName string) (
		port uint16, err error)
	KeepPortForward(ctx context.Context, port uint16, gateway netip.Addr,
		serverName string, logger Logger) (err error)
}

type NoPortForwarding struct {
	providerName string
}

func NewNoPortForwarding(providerName string) *NoPortForwarding {
	return &NoPortForwarding{
		providerName: providerName,
	}
}

var ErrPortForwardingNotSupported = errors.New("custom port forwarding obtention is not supported")

func (n *NoPortForwarding) PortForward(context.Context, *http.Client,
	Logger, netip.Addr, string) (port uint16, err error) {
	return 0, fmt.Errorf("%w: for %s", ErrPortForwardingNotSupported, n.providerName)
}

func (n *NoPortForwarding) KeepPortForward(context.Context, uint16, netip.Addr,
	string, Logger) (err error) {
	return fmt.Errorf("%w: for %s", ErrPortForwardingNotSupported, n.providerName)
}
