package resolver

import (
	"context"
	"net"
)

func newResolver(resolverAddress string) *net.Resolver {
	d := net.Dialer{}
	resolverAddress = net.JoinHostPort(resolverAddress, "53")
	return &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, _, _ string) (net.Conn, error) {
			return d.DialContext(ctx, "udp", resolverAddress)
		},
	}
}
