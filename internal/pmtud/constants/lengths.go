package constants

const (
	MaxEthernetFrameSize uint32 = 1500
	// MinIPv4MTU is defined according to
	// https://en.wikipedia.org/wiki/Maximum_transmission_unit#MTUs_for_common_media
	MinIPv4MTU uint32 = 68
	MinIPv6MTU uint32 = 1280

	IPv4HeaderLength uint32 = 20
	IPv6HeaderLength uint32 = 40
	UDPHeaderLength  uint32 = 8
	// BaseTCPHeaderLength is the TCP header length without options,
	// which is the minimum TCP header length.
	BaseTCPHeaderLength uint32 = 20
	// MaxTCPHeaderLength is the TCP header length with the maximum options length of 40 bytes.
	// Note this is a hard maximum because of the 4-bit data offset field in the TCP header (15x4=60).
	MaxTCPHeaderLength     uint32 = 60
	WireguardHeaderLength  uint32 = 32
	OpenVPNHeaderMaxLength uint32 = 1 + // opcode
		8 + // session id
		4 + // packet id
		28 // max possible auth tag/iv
)
