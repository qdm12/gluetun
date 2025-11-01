package portforward

import (
	"context"
	"net/netip"
	"os/exec"
)

type Service interface {
	Start(ctx context.Context) (runError <-chan error, err error)
	Stop() (err error)
	GetPortsForwarded() (ports []uint16)
	SetPortsForwarded(ctx context.Context, ports []uint16) (err error)
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

type Cmder interface {
	Start(cmd *exec.Cmd) (stdoutLines, stderrLines <-chan string,
		waitError <-chan error, startErr error)
}
