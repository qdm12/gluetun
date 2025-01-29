package pmtud

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	"fmt"
	"math/rand/v2"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

func buildMessageToSend(ipVersion string, mtu int) (id uint16, message *icmp.Message) {
	var seed [32]byte
	_, _ = cryptorand.Read(seed[:])
	randomSource := rand.NewChaCha8(seed)

	const uint16Bytes = 2
	idBytes := make([]byte, uint16Bytes)
	_, _ = randomSource.Read(idBytes)
	id = binary.BigEndian.Uint16(idBytes)

	var ipHeaderLength int
	var icmpType icmp.Type
	switch ipVersion {
	case "v4":
		ipHeaderLength = ipv4.HeaderLen
		icmpType = ipv4.ICMPTypeEcho
	case "v6":
		ipHeaderLength = ipv6.HeaderLen
		icmpType = ipv6.ICMPTypeEchoRequest
	default:
		panic(fmt.Sprintf("IP version %q not supported", ipVersion))
	}
	const pingHeaderLength = 0 +
		1 + // type
		1 + // code
		2 + // checksum
		2 + // identifier
		2 // sequence number
	pingBodyDataSize := mtu - ipHeaderLength - pingHeaderLength
	messageBodyData := make([]byte, pingBodyDataSize)
	_, _ = randomSource.Read(messageBodyData)

	// See https://www.iana.org/assignments/icmp-parameters/icmp-parameters.xhtml#icmp-parameters-types
	message = &icmp.Message{
		Type:     icmpType, // echo request
		Code:     0,        // no code
		Checksum: 0,        // calculated at encoding (ipv4) or sending (ipv6)
		Body: &icmp.Echo{
			ID:   int(id),
			Seq:  0, // only one packet
			Data: messageBodyData,
		},
	}
	return id, message
}
