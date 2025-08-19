package pmtud

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"net/netip"
	"runtime"
	"syscall"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

const (
	// see https://en.wikipedia.org/wiki/Maximum_transmission_unit#MTUs_for_common_media
	minIPv4MTU     = 68
	icmpv4Protocol = 1
)

func listenICMPv4(ctx context.Context) (conn net.PacketConn, err error) {
	var listenConfig net.ListenConfig
	listenConfig.Control = func(_, _ string, rawConn syscall.RawConn) error {
		var setDFErr error
		err := rawConn.Control(func(fd uintptr) {
			setDFErr = setDontFragment(fd) // runs when calling ListenPacket
		})
		if err == nil {
			err = setDFErr
		}
		return err
	}

	const listenAddress = ""
	packetConn, err := listenConfig.ListenPacket(ctx, "ip4:icmp", listenAddress)
	if err != nil {
		return nil, fmt.Errorf("listening for ICMP packets: %w", err)
	}

	if runtime.GOOS == "darwin" || runtime.GOOS == "ios" {
		packetConn = ipv4ToNetPacketConn(ipv4.NewPacketConn(packetConn))
	}

	return packetConn, nil
}

func findIPv4NextHopMTU(ctx context.Context, ip netip.Addr,
	physicalLinkMTU int, pingTimeout time.Duration, logger Logger,
) (mtu int, err error) {
	if ip.Is6() {
		panic("IP address is not v4")
	}
	conn, err := listenICMPv4(ctx)
	if err != nil {
		return 0, fmt.Errorf("listening for ICMP packets: %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()
	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	// First try to send a packet which is too big to get the maximum MTU
	// directly.
	outboundID, outboundMessage := buildMessageToSend("v4", physicalLinkMTU)
	encodedMessage, err := outboundMessage.Marshal(nil)
	if err != nil {
		return 0, fmt.Errorf("encoding ICMP message: %w", err)
	}

	_, err = conn.WriteTo(encodedMessage, &net.IPAddr{IP: ip.AsSlice()})
	if err != nil {
		return 0, fmt.Errorf("writing ICMP message: %w", err)
	}

	buffer := make([]byte, physicalLinkMTU)

	for { // for loop in case we read an echo reply for another ICMP request
		// Note we need to read the whole packet in one call to ReadFrom, so the buffer
		// must be large enough to read the entire reply packet. See:
		// https://groups.google.com/g/golang-nuts/c/5dy2Q4nPs08/m/KmuSQAGEtG4J
		bytesRead, _, err := conn.ReadFrom(buffer)
		if err != nil {
			return 0, fmt.Errorf("reading from ICMP connection: %w", err)
		}
		packetBytes := buffer[:bytesRead]
		// Side note: echo reply should be at most the number of bytes
		// sent, and can be lower, more precisely 576-ipHeader bytes,
		// in case the next hop we are reaching replies with a destination
		// unreachable and wants to ensure the response makes it way back
		// by keeping a low packet size, see:
		// https://datatracker.ietf.org/doc/html/rfc1122#page-59

		inboundMessage, err := icmp.ParseMessage(icmpv4Protocol, packetBytes)
		if err != nil {
			return 0, fmt.Errorf("parsing message: %w", err)
		}

		switch typedBody := inboundMessage.Body.(type) {
		case *icmp.DstUnreach:
			const fragmentationRequiredAndDFFlagSetCode = 4
			if inboundMessage.Code != fragmentationRequiredAndDFFlagSetCode {
				return 0, fmt.Errorf("%w: code %d",
					ErrICMPDestinationUnreachable, inboundMessage.Code)
			}

			// See https://datatracker.ietf.org/doc/html/rfc1191#section-4
			// Note: the go library does not handle this NextHopMTU section.
			nextHopMTU := packetBytes[6:8]
			mtu = int(binary.BigEndian.Uint16(nextHopMTU))
			err = checkMTU(mtu, minIPv4MTU, physicalLinkMTU)
			if err != nil {
				return 0, fmt.Errorf("checking next-hop-mtu found: %w", err)
			}

			// The code below is really for sanity checks
			packetBytes = packetBytes[8:]
			header, err := ipv4.ParseHeader(packetBytes)
			if err != nil {
				return 0, fmt.Errorf("parsing IPv4 header: %w", err)
			}
			packetBytes = packetBytes[header.Len:] // truncated original datagram

			const truncated = true
			err = checkEchoReply(icmpv4Protocol, packetBytes, outboundMessage, truncated)
			if err != nil {
				return 0, fmt.Errorf("checking echo reply: %w", err)
			}
			return mtu, nil
		case *icmp.Echo:
			inboundID := uint16(typedBody.ID) //nolint:gosec
			if inboundID == outboundID {
				return physicalLinkMTU, nil
			}
			logger.Debug("discarding received ICMP echo reply with id %d mismatching sent id %d",
				inboundID, outboundID)
			continue
		default:
			return 0, fmt.Errorf("%w: %T", ErrICMPBodyUnsupported, typedBody)
		}
	}
}
