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

// PathMTUDiscover discovers the maximum MTU using both ICMP and TCP.
// Multiple ICMP addresses and TCP addresses can be specified for redundancy.
// ICMP PMTUD is run first, then TCP PMTUD is run with the maximum MTU found from
// ICMP PMTUD as its upper bound.
// If the physicalLinkMTU is zero, it defaults to 1500 which is the ethernet standard MTU.
// If the pingTimeout is zero, it defaults to 1 second.
// If the logger is nil, a no-op logger is used.
// It returns [ErrMTUNotFound] if the MTU could not be determined.
func PathMTUDiscover(ctx context.Context, icmpAddrs []netip.Addr, tcpAddrs []netip.AddrPort,
	physicalLinkMTU uint32, tryTimeout time.Duration, logger Logger) (
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
	maxPossibleMTU := physicalLinkMTU
	for _, icmpIP := range icmpAddrs {
		mtu, err := icmp.PathMTUDiscover(ctx, icmpIP, physicalLinkMTU,
			tryTimeout, logger)
		switch {
		case err == nil:
			logger.Debugf("ICMP path MTU discovery against %s found maximum valid MTU %d", icmpIP, mtu)
			maxPossibleMTU = mtu
		case errors.Is(err, icmp.ErrNotPermitted), errors.Is(err, icmp.ErrMTUNotFound):
			logger.Debugf("ICMP path MTU discovery failed: %s", err)
		default:
			return 0, fmt.Errorf("ICMP path MTU discovery: %w", err)
		}
	}

	for _, addrPort := range tcpAddrs {
		mtu, err = tcp.PathMTUDiscover(ctx, addrPort, maxPossibleMTU, logger)
		if err != nil {
			logger.Debugf("TCP path MTU discovery to %s failed: %s", addrPort, err)
			continue
		}
		logger.Debugf("TCP path MTU discovery to %s found maximum valid MTU %d", addrPort, mtu)
		return mtu, nil
	}
	return 0, fmt.Errorf("TCP path MTU discovery: last error: %w", err)
}
