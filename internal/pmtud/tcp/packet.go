package tcp

import (
	"encoding/binary"
	"math/rand/v2"
	"net/netip"

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
		// Pseudo-header, this is actually not part of the packet since
		// the kernel will calculate and add it itself to the packet;
		// it is only used for calculating the TCP checksum.
		ipHeader = ip.HeaderV6(src.Addr(), dst.Addr(),
			uint16(payloadLength), byte(constants.IPPROTO_TCP)) //nolint:gosec
	}

	tcpHeader := makeTCPHeader(src.Port(), dst.Port(), seq, ack, flags)

	dataLength := int(payloadLength - constants.BaseTCPHeaderLength)
	var data []byte
	if dataLength > 0 {
		data = generatePayload(uint16(dataLength)) //nolint:gosec
	}
	checksum := tcpChecksum(ipHeader, tcpHeader, data)
	tcpHeader[16] = byte(checksum >> 8)   //nolint:mnd
	tcpHeader[17] = byte(checksum & 0xff) //nolint:mnd

	var packet []byte
	i := 0
	if dst.Addr().Is4() {
		packet = make([]byte, len(ipHeader)+int(constants.BaseTCPHeaderLength)+dataLength)
		copy(packet, ipHeader)
		i += len(ipHeader)
	} else {
		packet = make([]byte, int(constants.BaseTCPHeaderLength)+dataLength)
	}
	copy(packet[i:], tcpHeader)
	i += int(constants.BaseTCPHeaderLength)
	copy(packet[i:], data)
	return packet
}

// generatePayload creates a byte slice of 'length' size.
// For lengths below 88B, it returns pseudo random data.
// For lengths above, it returns a structured TLS Client Hello with padding,
// which is more likely to be accepted by servers and not trigger RST replies.
//
//nolint:mnd
func generatePayload(length uint16) []byte {
	const minTLSClientHelloSize = 5 + // TLS record
		4 + // handshake header
		67 + // client hello
		4 + // cipher suites
		2 + // compression methods
		2 + // extensions length
		4 // padding extension header
	if length < minTLSClientHelloSize {
		data := make([]byte, length)
		makeRandom(data)
		return data
	}

	payload := make([]byte, length)

	// --- TLS Record Layer ---
	payload[0] = 0x16 // Handshake
	payload[1] = 0x03 // Version 3.1
	payload[2] = 0x01
	binary.BigEndian.PutUint16(payload[3:5], length-5)

	// --- Handshake Header ---
	payload[5] = 0x01 // Client Hello
	handshakeLength := make([]byte, 4)
	// TLS Handshake length is 24-bit.
	// We use a 4-byte buffer and copy the trailing 3 bytes.
	binary.BigEndian.PutUint32(handshakeLength, uint32(length-9))
	copy(payload[6:9], handshakeLength[1:])

	// --- Client Hello Body ---
	payload[9] = 0x03 // Version 3.3 (TLS 1.2)
	payload[10] = 0x03
	makeRandom(payload[11:43]) // 32 bytes of random
	payload[43] = 32           // Session ID length

	// Cipher Suites (Length: 2, Data: 2)
	binary.BigEndian.PutUint16(payload[44:46], 2)
	binary.BigEndian.PutUint16(payload[46:48], 0x009c) // TLS_RSA_WITH_AES_128_GCM_SHA256

	payload[48] = 0x01 // Compression length
	payload[49] = 0x00 // Null compression

	// --- Extensions ---
	binary.BigEndian.PutUint16(payload[50:52], length-52) // extension length

	// --- Padding Extension (Type 21) ---
	binary.BigEndian.PutUint16(payload[52:54], 21)
	const bytesUsedSoFar = 88
	paddingDataLength := length - bytesUsedSoFar
	binary.BigEndian.PutUint16(payload[54:56], paddingDataLength)

	return payload
}

func makeRandom(b []byte) {
	for i := range b {
		b[i] = byte(rand.Uint32()) //nolint:gosec
	}
}
