package utils

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
)

type NoPortForwarder interface {
	PortForward(ctx context.Context, client *http.Client,
		logger Logger, gateway net.IP, serverName string) (
		port uint16, err error)
	KeepPortForward(ctx context.Context, gateway net.IP,
		serverName string) (err error)
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
	Logger, net.IP, string) (port uint16, err error) {
	return 0, fmt.Errorf("%w: for %s", ErrPortForwardingNotSupported, n.providerName)
}

func (n *NoPortForwarding) KeepPortForward(context.Context, net.IP, string) (err error) {
	return fmt.Errorf("%w: for %s", ErrPortForwardingNotSupported, n.providerName)
}
