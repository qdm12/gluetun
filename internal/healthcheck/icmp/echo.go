package icmp

import (
	"bytes"
	"context"
	cryptorand "crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net"
	"net/netip"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

var (
	ErrICMPBodyUnsupported  = errors.New("ICMP body type is not supported")
	ErrICMPDstUnreachable   = errors.New("ICMP destination unreachable")
	ErrICMPEchoDataMismatch = errors.New("ICMP data mismatch")
)

type Echoer struct {
	buffer       []byte
	randomSource io.Reader
	warner       Warner
}

func NewEchoer(warner Warner) *Echoer {
	const maxICMPEchoSize = 1500
	buffer := make([]byte, maxICMPEchoSize)
	var seed [32]byte
	_, _ = cryptorand.Read(seed[:])
	randomSource := rand.NewChaCha8(seed)
	return &Echoer{
		buffer:       buffer,
		randomSource: randomSource,
		warner:       warner,
	}
}

var ErrTimedOut = errors.New("timed out waiting for ICMP echo reply")

func (i *Echoer) Echo(ctx context.Context, ip netip.Addr) (err error) {
	var ipVersion string
	var conn net.PacketConn
	if ip.Is4() {
		ipVersion = "v4"
		conn, err = listenICMPv4(ctx)
	} else {
		ipVersion = "v6"
		conn, err = listenICMPv6(ctx)
	}
	if err != nil {
		return fmt.Errorf("listening for ICMP packets: %w", err)
	}

	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	const echoDataSize = 32
	id, message := buildMessageToSend(ipVersion, echoDataSize, i.randomSource)

	encodedMessage, err := message.Marshal(nil)
	if err != nil {
		return fmt.Errorf("encoding ICMP message: %w", err)
	}

	_, err = conn.WriteTo(encodedMessage, &net.IPAddr{IP: ip.AsSlice()})
	if err != nil {
		return fmt.Errorf("writing ICMP message: %w", err)
	}

	receivedData, err := receiveEchoReply(conn, id, i.buffer, ipVersion, i.warner)
	if err != nil {
		if errors.Is(err, net.ErrClosed) && ctx.Err() != nil {
			return fmt.Errorf("%w", ErrTimedOut)
		}
		return fmt.Errorf("receiving ICMP echo reply: %w", err)
	}

	sentData := message.Body.(*icmp.Echo).Data //nolint:forcetypeassert
	if !bytes.Equal(receivedData, sentData) {
		return fmt.Errorf("%w: sent %x and received %x", ErrICMPEchoDataMismatch, sentData, receivedData)
	}

	return nil
}

func buildMessageToSend(ipVersion string, size uint, randomSource io.Reader) (id int, message *icmp.Message) {
	const uint16Bytes = 2
	idBytes := make([]byte, uint16Bytes)
	_, _ = randomSource.Read(idBytes)
	id = int(binary.BigEndian.Uint16(idBytes))

	var icmpType icmp.Type
	switch ipVersion {
	case "v4":
		icmpType = ipv4.ICMPTypeEcho
	case "v6":
		icmpType = ipv6.ICMPTypeEchoRequest
	default:
		panic(fmt.Sprintf("IP version %q not supported", ipVersion))
	}
	messageBodyData := make([]byte, size)
	_, _ = randomSource.Read(messageBodyData)

	// See https://www.iana.org/assignments/icmp-parameters/icmp-parameters.xhtml#icmp-parameters-types
	message = &icmp.Message{
		Type:     icmpType, // echo request
		Code:     0,        // no code
		Checksum: 0,        // calculated at encoding (ipv4) or sending (ipv6)
		Body: &icmp.Echo{
			ID:   id,
			Seq:  0, // only one packet
			Data: messageBodyData,
		},
	}
	return id, message
}

func receiveEchoReply(conn net.PacketConn, id int, buffer []byte, ipVersion string, logger Warner,
) (data []byte, err error) {
	var icmpProtocol int
	const (
		icmpv4Protocol = 1
		icmpv6Protocol = 58
	)
	switch ipVersion {
	case "v4":
		icmpProtocol = icmpv4Protocol
	case "v6":
		icmpProtocol = icmpv6Protocol
	default:
		panic(fmt.Sprintf("unknown IP version: %s", ipVersion))
	}

	for {
		// Note we need to read the whole packet in one call to ReadFrom, so the buffer
		// must be large enough to read the entire reply packet. See:
		// https://groups.google.com/g/golang-nuts/c/5dy2Q4nPs08/m/KmuSQAGEtG4J
		bytesRead, _, err := conn.ReadFrom(buffer)
		if err != nil {
			return nil, fmt.Errorf("reading from ICMP connection: %w", err)
		}
		packetBytes := buffer[:bytesRead]

		// Parse the ICMP message
		message, err := icmp.ParseMessage(icmpProtocol, packetBytes)
		if err != nil {
			return nil, fmt.Errorf("parsing message: %w", err)
		}

		switch body := message.Body.(type) {
		case *icmp.Echo:
			if id != body.ID {
				logger.Warnf("ignoring ICMP reply mismatching expected id %d (id: %d, type: %d, code: %d, length: %d)",
					id, body.ID, message.Type, message.Code, len(packetBytes))
				continue // not the ID we are looking for
			}
			return body.Data, nil
		case *icmp.DstUnreach:
			return nil, fmt.Errorf("%w (id %d and reply ICMP type 3 code %d)", ErrICMPDstUnreachable, id, message.Code)
		default:
			return nil, fmt.Errorf("%w: %T (id %d, type %d)", ErrICMPBodyUnsupported, body, id, message.Type)
		}
	}
}
