package service

import (
	"context"
	"net/netip"

	"github.com/qdm12/gluetun/internal/provider/utils"
)

type PortAllower interface {
	SetAllowedPort(ctx context.Context, port uint16, intf string) (err error)
	RemoveAllowedPort(ctx context.Context, port uint16) (err error)
}

type Routing interface {
	VPNLocalGatewayIP(vpnInterface string) (gateway netip.Addr, err error)
}

type Logger interface {
	Debug(s string)
	Info(s string)
	Warn(s string)
	Error(s string)
}

type PortForwarder interface {
	Name() string
	PortForward(ctx context.Context, objects utils.PortForwardObjects) (
		port uint16, err error)
	KeepPortForward(ctx context.Context, objects utils.PortForwardObjects) (err error)
}
