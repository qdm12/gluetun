package icmp

import (
	"context"
	"errors"
	"fmt"
	"net/netip"
	"time"

	"github.com/qdm12/gluetun/internal/pmtud/constants"
)

// PathMTUDiscover discovers the path MTU to the given IP address
// using ICMP.
// It first tries to get the next hop MTU using ICMP messages.
// If that fails, it falls back to sending echo requests with
// different packet sizes to find the maximum MTU.
// The function returns [ErrMTUNotFound] if the MTU could not be determined.
func PathMTUDiscover(ctx context.Context, ip netip.Addr,
	physicalLinkMTU uint32, timeout time.Duration, logger Logger,
) (mtu uint32, err error) {
	if ip.Is4() {
		logger.Debug("finding IPv4 next hop MTU")
		mtu, err = findIPv4NextHopMTU(ctx, ip, physicalLinkMTU, timeout, logger)
		switch {
		case err == nil:
			return mtu, nil
		case errors.Is(err, errTimeout) || errors.Is(err, ErrCommunicationAdministrativelyProhibited): // blackhole
		default:
			return 0, fmt.Errorf("finding IPv4 next hop MTU: %w", err)
		}
	} else {
		logger.Debug("requesting IPv6 ICMP packet-too-big reply")
		mtu, err = getIPv6PacketTooBig(ctx, ip, physicalLinkMTU, timeout, logger)
		switch {
		case err == nil:
			return mtu, nil
		case errors.Is(err, errTimeout): // blackhole
		default:
			return 0, fmt.Errorf("getting IPv6 packet-too-big message: %w", err)
		}
	}

	// Fall back method: send echo requests with different packet
	// sizes and check which ones succeed to find the maximum MTU.
	logger.Debug("falling back to sending different sized echo packets")
	minMTU := constants.MinIPv4MTU
	if ip.Is6() {
		minMTU = constants.MinIPv6MTU
	}
	return pmtudMultiSizes(ctx, ip, minMTU, physicalLinkMTU, timeout, logger)
}
