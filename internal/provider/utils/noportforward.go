package utils

import (
	"context"
	"errors"
	"fmt"
)

type NoPortForwarder interface {
	PortForward(ctx context.Context, objects PortForwardObjects) (port uint16, err error)
	KeepPortForward(ctx context.Context, objects PortForwardObjects) (err error)
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

func (n *NoPortForwarding) PortForward(context.Context, PortForwardObjects) (
	port uint16, err error) {
	return 0, fmt.Errorf("%w: for %s", ErrPortForwardingNotSupported, n.providerName)
}

func (n *NoPortForwarding) KeepPortForward(context.Context, PortForwardObjects) (err error) {
	return fmt.Errorf("%w: for %s", ErrPortForwardingNotSupported, n.providerName)
}
