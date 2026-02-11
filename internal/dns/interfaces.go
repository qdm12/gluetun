package dns

import (
	"context"
	"net/netip"
)

type Logger interface {
	Debug(s string)
	Info(s string)
	Warn(s string)
	Error(s string)
}

type Firewall interface {
	RestrictOutputAddrPort(ctx context.Context, addrPort netip.AddrPort) (err error)
}
