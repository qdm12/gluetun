package utils

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/qdm12/golibs/logging"
)

type NoPortForwarder interface {
	PortForward(ctx context.Context, client *http.Client,
		logger logging.Logger, gateway net.IP, serverName string) (
		port uint16, err error)
	KeepPortForward(ctx context.Context, client *http.Client,
		logger logging.Logger, port uint16, gateway net.IP, serverName string) (
		err error)
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

func (n *NoPortForwarding) PortForward(ctx context.Context, client *http.Client,
	logger logging.Logger, gateway net.IP, serverName string) (
	port uint16, err error) {
	return 0, fmt.Errorf("%w: for %s", ErrPortForwardingNotSupported, n.providerName)
}

func (n *NoPortForwarding) KeepPortForward(ctx context.Context, client *http.Client,
	logger logging.Logger, port uint16, gateway net.IP, serverName string) (
	err error) {
	return fmt.Errorf("%w: for %s", ErrPortForwardingNotSupported, n.providerName)
}
