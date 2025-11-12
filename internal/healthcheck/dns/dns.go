package dns

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/netip"

	"github.com/qdm12/dns/v2/pkg/provider"
)

// Client is a simple plaintext UDP DNS client, to be used for healthchecks.
// Note the client connects to a DNS server only over UDP on port 53,
// because we don't want to use DoT or DoH and impact the TCP connections
// when running a healthcheck.
type Client struct {
	serverAddrs []netip.AddrPort
	dnsIPIndex  int
}

func New() *Client {
	return &Client{
		serverAddrs: concatAddrPorts([][]netip.AddrPort{
			provider.Cloudflare().Plain.IPv4,
			provider.Google().Plain.IPv4,
			provider.Quad9().Plain.IPv4,
			provider.OpenDNS().Plain.IPv4,
			provider.LibreDNS().Plain.IPv4,
			provider.Quadrant().Plain.IPv4,
			provider.CiraProtected().Plain.IPv4,
		}),
	}
}

func concatAddrPorts(addrs [][]netip.AddrPort) []netip.AddrPort {
	var result []netip.AddrPort
	for _, addrList := range addrs {
		result = append(result, addrList...)
	}
	return result
}

var ErrLookupNoIPs = errors.New("no IPs found from DNS lookup")

func (c *Client) Check(ctx context.Context) error {
	dnsAddr := c.serverAddrs[c.dnsIPIndex].String()
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, _, _ string) (net.Conn, error) {
			dialer := net.Dialer{}
			return dialer.DialContext(ctx, "udp", dnsAddr)
		},
	}
	ips, err := resolver.LookupIP(ctx, "ip", "github.com")
	switch {
	case err != nil:
		c.dnsIPIndex = (c.dnsIPIndex + 1) % len(c.serverAddrs)
		return err
	case len(ips) == 0:
		c.dnsIPIndex = (c.dnsIPIndex + 1) % len(c.serverAddrs)
		return fmt.Errorf("%w", ErrLookupNoIPs)
	default:
		return nil
	}
}
