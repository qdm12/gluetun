package tcp

import (
	"math/rand/v2"
	"net/netip"
	"syscall"

	"github.com/qdm12/gluetun/internal/pmtud/constants"
	"github.com/qdm12/gluetun/internal/pmtud/ip"
)

// createSYNPacket creates a TCP SYN packet for initiating a handshake.
// SYN packets have normally no data payload, so you SHOULD set mtu to 0.
// However, in some cases where the server closes the connection with RST immediately,
// it can be useful to add some data payload to a SYN packet and check if the server still
// replies. Only set mtu to a non zero value if you know what you are doing.
func createSYNPacket(src, dst netip.AddrPort, mtu uint32) (packet []byte, seq uint32) {
	seq = rand.Uint32()                            //nolint:gosec
	const ack = 0                                  // SYN has no ACK number
	payloadLength := constants.BaseTCPHeaderLength // no data payload
	if mtu > 0 {
		payloadLength = getPayloadLength(mtu, dst)
	}
	return createPacket(src, dst, seq, ack, payloadLength, synFlag), seq
}

// createACKPacket creates a TCP ACK packet.
// If the mtu is set to 0, no payload is sent.
// Otherwise, the payload is calculated to test the MTU given.
func createACKPacket(src, dst netip.AddrPort, seq, ack uint32, mtu uint32) []byte {
	payloadLength := constants.BaseTCPHeaderLength // no data payload
	if mtu > 0 {
		payloadLength = getPayloadLength(mtu, dst)
	}
	const flags = ackFlag | pshFlag
	return createPacket(src, dst, seq, ack, payloadLength, flags)
}

func createRSTPacket(src, dst netip.AddrPort, seq, ack uint32) []byte {
	const payloadLength = constants.BaseTCPHeaderLength // no data payload
	return createPacket(src, dst, seq, ack, payloadLength, rstFlag)
}

func getPayloadLength(mtu uint32, dst netip.AddrPort) uint32 {
	var ipHeaderLength uint32
	if dst.Addr().Is4() {
		ipHeaderLength = constants.IPv4HeaderLength
	} else {
		ipHeaderLength = constants.IPv6HeaderLength
	}
	if mtu < ipHeaderLength+constants.BaseTCPHeaderLength {
		panic("MTU too small to hold IP and TCP headers")
	}
	return mtu - ipHeaderLength
}

func createPacket(src, dst netip.AddrPort,
	seq, ack, payloadLength uint32, flags byte,
) []byte {
	if payloadLength < constants.BaseTCPHeaderLength {
		panic("payload length is too small to hold TCP header")
	}

	var ipHeader []byte
	if dst.Addr().Is4() {
		ipHeader = ip.HeaderV4(src.Addr(), dst.Addr(), payloadLength)
	} else {
		ipHeader = ip.HeaderV6(src.Addr(), dst.Addr(),
			uint16(payloadLength), byte(syscall.IPPROTO_TCP)) //nolint:gosec
	}

	tcpHeader := makeTCPHeader(src.Port(), dst.Port(), seq, ack, flags)

	// data is just zeroes
	dataLength := int(payloadLength) - int(constants.BaseTCPHeaderLength)
	var data []byte
	if dataLength > 0 {
		data = make([]byte, dataLength)
	}
	checksum := tcpChecksum(ipHeader, tcpHeader, data)
	tcpHeader[16] = byte(checksum >> 8)   //nolint:mnd
	tcpHeader[17] = byte(checksum & 0xff) //nolint:mnd

	packet := make([]byte, len(ipHeader)+int(constants.BaseTCPHeaderLength)+dataLength)
	copy(packet, ipHeader)
	copy(packet[len(ipHeader):], tcpHeader)
	copy(packet[len(ipHeader)+int(constants.BaseTCPHeaderLength):], data)
	return packet
}
