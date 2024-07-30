package portforward

import (
	"context"
	"net/netip"
)

type Service interface {
	Start(ctx context.Context) (runError <-chan error, err error)
	Stop() (err error)
	GetPortsForwarded() (ports []uint16)
}

type Routing interface {
	VPNLocalGatewayIP(vpnInterface string) (gateway netip.Addr, err error)
	AssignedIP(interfaceName string, family int) (ip netip.Addr, err error)
}

type PortAllower interface {
	SetAllowedPort(ctx context.Context, port uint16, intf string) (err error)
	RemoveAllowedPort(ctx context.Context, port uint16) (err error)
	RedirectPort(ctx context.Context, intf string, sourcePort,
		destinationPort uint16) (err error)
}

type Logger interface {
	Debug(s string)
	Info(s string)
	Warn(s string)
	Error(s string)
}
