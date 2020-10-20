package updater

import (
	"context"
	"net"
)

type (
	lookupIPFunc func(ctx context.Context, host string) (ips []net.IP, err error)
)
