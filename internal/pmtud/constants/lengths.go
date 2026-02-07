package constants

const (
	MaxEthernetFrameSize uint32 = 1500
	// MinIPv4MTU is defined according to
	// https://en.wikipedia.org/wiki/Maximum_transmission_unit#MTUs_for_common_media
	MinIPv4MTU uint32 = 68
	MinIPv6MTU uint32 = 1280

	IPv4HeaderLength       uint32 = 20
	IPv6HeaderLength       uint32 = 40
	UDPHeaderLength        uint32 = 8
	TCPHeaderLength        uint32 = 20
	WireguardHeaderLength  uint32 = 32
	OpenVPNHeaderMaxLength uint32 = 1 + // opcode
		8 + // session id
		4 + // packet id
		28 // max possible auth tag/iv
)
