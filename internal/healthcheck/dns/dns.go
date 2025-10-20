package dns

import (
	"context"
	"errors"
	"fmt"
	"net"
)

// Client is a simple plaintext UDP DNS client, to be used for healthchecks.
// Note the client connects to a DNS server only over UDP on port 53,
// because we don't want to use DoT or DoH and impact the TCP connections
// when running a healthcheck.
type Client struct{}

func New() *Client {
	return &Client{}
}

var ErrLookupNoIPs = errors.New("no IPs found from DNS lookup")

func (c *Client) Check(ctx context.Context) error {
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, _, _ string) (net.Conn, error) {
			dialer := net.Dialer{}
			return dialer.DialContext(ctx, "udp", "1.1.1.1:53")
		},
	}
	ips, err := resolver.LookupIP(ctx, "ip", "github.com")
	switch {
	case err != nil:
		return err
	case len(ips) == 0:
		return fmt.Errorf("%w", ErrLookupNoIPs)
	default:
		return nil
	}
}
