//go:build !darwin

package tcp

import (
	"github.com/qdm12/gluetun/internal/pmtud/constants"
)

func stripIPv4Header(reply []byte) (result []byte, ok bool) {
	if len(reply) < int(constants.IPv4HeaderLength) {
		return nil, false // not an IPv4 packet
	}

	version := reply[0] >> 4 //nolint:mnd
	const ipv4Version = 4
	if version != ipv4Version {
		return nil, false
	}
	// For IPv4 we need to skip the IP header, which is at least
	// 20B and can be up to 60B.
	// The Internet Header Length is the lower 4 bits of the first byte and
	// represents the number of 32-bit words of the header length.
	const ihlMask byte = 0x0F
	const bytesInWord = 4
	headerLength := int((reply[0] & ihlMask)) * bytesInWord
	if len(reply) < headerLength {
		return nil, false // not enough data for full IPv4 header
	}
	return reply[headerLength:], true
}
