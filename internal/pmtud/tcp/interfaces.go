package tcp

import (
	"context"
	"net/netip"
)

type Firewall interface {
	TempDropOutputTCPRST(ctx context.Context, addrPort netip.AddrPort,
		excludeMark int) (revert func(ctx context.Context) error, err error)
}

type Logger interface {
	Debug(msg string)
	Debugf(msg string, args ...any)
	Warnf(msg string, args ...any)
}
