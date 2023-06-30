package natpmp

import (
	"context"
	"encoding/binary"
	"fmt"
	"net/netip"
	"time"
)

// ExternalAddress fetches the duration since the start of epoch and the external
// IPv4 address of the gateway.
// See https://www.ietf.org/rfc/rfc6886.html#section-3.2
func (c *Client) ExternalAddress(ctx context.Context, gateway netip.Addr) (
	durationSinceStartOfEpoch time.Duration,
	externalIPv4Address netip.Addr, err error) {
	request := []byte{0, 0} // version 0, operationCode 0
	const responseSize = 12
	response, err := c.rpc(ctx, gateway, request, responseSize)
	if err != nil {
		return 0, externalIPv4Address, fmt.Errorf("executing remote procedure call: %w", err)
	}

	secondsSinceStartOfEpoch := binary.BigEndian.Uint32(response[4:8])
	durationSinceStartOfEpoch = time.Duration(secondsSinceStartOfEpoch) * time.Second
	externalIPv4Address = netip.AddrFrom4([4]byte{response[8], response[9], response[10], response[11]})
	return durationSinceStartOfEpoch, externalIPv4Address, nil
}
