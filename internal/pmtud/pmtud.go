package pmtud

import (
	"context"
	"errors"
	"fmt"
	"net/netip"
	"time"

	"github.com/qdm12/gluetun/internal/pmtud/icmp"
	"github.com/qdm12/gluetun/internal/pmtud/tcp"
)

// PathMTUDiscover discovers the maximum MTU for the path to the given ip address.
// If the physicalLinkMTU is zero, it defaults to 1500 which is the ethernet standard MTU.
// If the pingTimeout is zero, it defaults to 1 second.
// If the logger is nil, a no-op logger is used.
// It returns [ErrMTUNotFound] if the MTU could not be determined.
func PathMTUDiscover(ctx context.Context, icmpIP netip.Addr,
	tcpAddresses []netip.AddrPort, physicalLinkMTU uint32,
	tryTimeout time.Duration, logger Logger) (
	mtu uint32, err error,
) {
	if physicalLinkMTU == 0 {
		const ethernetStandardMTU = 1500
		physicalLinkMTU = ethernetStandardMTU
	}
	if tryTimeout == 0 {
		tryTimeout = time.Second
	}
	if logger == nil {
		logger = &noopLogger{}
	}

	// Try finding the MTU using ICMP
	mtu, err = icmp.PathMTUDiscover(ctx, icmpIP, physicalLinkMTU,
		tryTimeout, logger)
	switch {
	case err == nil:
		logger.Debugf("ICMP path MTU discovery against %s found maximum valid MTU %d", icmpIP, mtu)
		return mtu, nil
	case errors.Is(err, icmp.ErrNotPermitted), errors.Is(err, icmp.ErrMTUNotFound):
		logger.Debugf("ICMP path MTU discovery failed: %s", err)
	default:
		return 0, fmt.Errorf("ICMP path MTU discovery: %w", err)
	}

	// If ICMP path MTU discovery is not permitted or failed completely,
	// we run the below TCP path MTU discovery.
	for _, tcpAddress := range tcpAddresses {
		mtu, err = tcp.PathMTUDiscover(ctx, tcpAddress, physicalLinkMTU, logger)
		if err != nil {
			logger.Debugf("TCP path MTU discovery to %s failed: %s", tcpAddress, err)
			continue
		}
		logger.Debugf("TCP path MTU discovery to %s found maximum valid MTU %d", tcpAddress, mtu)
		return mtu, nil
	}
	return 0, fmt.Errorf("TCP path MTU discovery: last error: %w", err)
}
