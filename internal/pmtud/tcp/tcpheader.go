package tcp

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/qdm12/gluetun/internal/pmtud/constants"
)

// For SYN, ack is 0.
// For SYN-ACK, ack is the sequence number + 1 of the SYN.
func makeTCPHeader(srcPort, dstPort uint16, seq, ack uint32, flags byte) []byte {
	header := make([]byte, constants.BaseTCPHeaderLength)
	binary.BigEndian.PutUint16(header[0:], srcPort)
	binary.BigEndian.PutUint16(header[2:], dstPort)
	binary.BigEndian.PutUint32(header[4:], seq)
	binary.BigEndian.PutUint32(header[8:], ack)
	//nolint:mnd
	header[12] = byte(constants.BaseTCPHeaderLength) << 2 // data offset
	header[13] = flags
	// windowSize can be left to 5840 even for IPv6, it doesn't matter.
	const windowSize = 5840
	binary.BigEndian.PutUint16(header[14:], windowSize)
	// header[16:17] is the checksum, set later
	// header[18:19] is urgent pointer, not needed for our use case
	return header
}

//nolint:mnd
func tcpChecksum(ipHeader, tcpHeader, payload []byte) uint16 {
	var pseudoHeader []byte
	isIPv6 := len(ipHeader) >= 40 && (ipHeader[0]>>4) == 6
	if isIPv6 {
		pseudoHeader = make([]byte, 40)
		copy(pseudoHeader[0:16], ipHeader[8:24])             // Source Address
		copy(pseudoHeader[16:32], ipHeader[24:40])           // Destination Address
		totalLength := uint32(len(tcpHeader) + len(payload)) //nolint:gosec
		binary.BigEndian.PutUint32(pseudoHeader[32:], totalLength)
		pseudoHeader[39] = 6 // Next Header (TCP)
	} else {
		pseudoHeader = make([]byte, 12)
		copy(pseudoHeader[0:4], ipHeader[12:16])
		copy(pseudoHeader[4:8], ipHeader[16:20])
		pseudoHeader[9] = 6
		totalLength := uint16(len(tcpHeader) + len(payload)) //nolint:gosec
		binary.BigEndian.PutUint16(pseudoHeader[10:], totalLength)
	}

	sum := uint32(0)
	for _, slice := range [][]byte{pseudoHeader, tcpHeader, payload} {
		for i := 0; i < len(slice)-1; i += 2 {
			sum += uint32(binary.BigEndian.Uint16(slice[i : i+2]))
		}
		if len(slice)%2 != 0 {
			sum += uint32(slice[len(slice)-1]) << 8
		}
	}
	for (sum >> 16) > 0 {
		sum = (sum & 0xFFFF) + (sum >> 16)
	}
	return ^uint16(sum) //nolint:gosec
}

const (
	tcpFlagsOffset      = 13
	rstFlag        byte = 0x04
	synFlag        byte = 0x02
	ackFlag        byte = 0x10
	pshFlag        byte = 0x08
)

type packetType uint8

const (
	packetTypeSYN packetType = iota + 1
	packetTypeSYNACK
	packetTypeACK
	packetTypeRST
)

func (p packetType) String() string {
	switch p {
	case packetTypeSYN:
		return "SYN"
	case packetTypeSYNACK:
		return "SYN-ACK"
	case packetTypeACK:
		return "ACK"
	case packetTypeRST:
		return "RST"
	default:
		panic("unknown packet type")
	}
}

var (
	errTCPHeaderTooShort    = errors.New("TCP header is too short")
	errTCPPacketTypeUnknown = errors.New("TCP packet type is unknown")
)

// parseTCPHeader parses some elements from the TCP header.
func parseTCPHeader(header []byte) (packetType packetType, seq, ack uint32, err error) {
	if len(header) < int(constants.BaseTCPHeaderLength) {
		return 0, 0, 0, fmt.Errorf("%w: %d bytes", errTCPHeaderTooShort, len(header))
	}
	flags := header[tcpFlagsOffset]
	switch {
	case (flags&synFlag) != 0 && (flags&ackFlag) == 0:
		packetType = packetTypeSYN
	case (flags&synFlag) != 0 && (flags&ackFlag) != 0:
		packetType = packetTypeSYNACK
	case (flags & rstFlag) != 0:
		packetType = packetTypeRST
	case (flags & ackFlag) != 0:
		packetType = packetTypeACK
	default:
		return 0, 0, 0, fmt.Errorf("%w: flags are 0x%02x", errTCPPacketTypeUnknown, flags)
	}

	seq = binary.BigEndian.Uint32(header[4:8])
	ack = binary.BigEndian.Uint32(header[8:12])
	return packetType, seq, ack, nil
}
