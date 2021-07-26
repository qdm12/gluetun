package utils

import (
	"context"
	"net"
	"net/http"

	"github.com/qdm12/gluetun/internal/firewall"
	"github.com/qdm12/golibs/logging"
)

type NoPortForwarder interface {
	PortForward(ctx context.Context, client *http.Client,
		pfLogger logging.Logger, gateway net.IP, portAllower firewall.PortAllower,
		syncState func(port uint16) (pfFilepath string))
}

type NoPortForwarding struct {
	providerName string
}

func NewNoPortForwarding(providerName string) *NoPortForwarding {
	return &NoPortForwarding{
		providerName: providerName,
	}
}

func (n *NoPortForwarding) PortForward(ctx context.Context, client *http.Client,
	pfLogger logging.Logger, gateway net.IP, portAllower firewall.PortAllower,
	syncState func(port uint16) (pfFilepath string)) {
	panic("custom port forwarding obtention is not supported for " + n.providerName)
}
