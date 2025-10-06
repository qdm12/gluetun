package pmtud

import (
	"context"
	"fmt"
	"net"
	"net/netip"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv6"
)

const (
	minIPv6MTU     = 1280
	icmpv6Protocol = 58
)

func listenICMPv6(ctx context.Context) (conn net.PacketConn, err error) {
	var listenConfig net.ListenConfig
	const listenAddress = ""
	packetConn, err := listenConfig.ListenPacket(ctx, "ip6:ipv6-icmp", listenAddress)
	if err != nil {
		return nil, fmt.Errorf("listening for ICMPv6 packets: %w", err)
	}
	return packetConn, nil
}

func getIPv6PacketTooBig(ctx context.Context, ip netip.Addr,
	physicalLinkMTU int, pingTimeout time.Duration, logger Logger,
) (mtu int, err error) {
	if ip.Is4() {
		panic("IP address is not v6")
	}
	conn, err := listenICMPv6(ctx)
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
	outboundID, outboundMessage := buildMessageToSend("v6", physicalLinkMTU)
	encodedMessage, err := outboundMessage.Marshal(nil)
	if err != nil {
		return 0, fmt.Errorf("encoding ICMP message: %w", err)
	}

	_, err = conn.WriteTo(encodedMessage, &net.IPAddr{IP: ip.AsSlice(), Zone: ip.Zone()})
	if err != nil {
		return 0, fmt.Errorf("writing ICMP message: %w", err)
	}

	buffer := make([]byte, physicalLinkMTU)

	for { // for loop if we encounter another ICMP packet with an unknown id.
		// Note we need to read the whole packet in one call to ReadFrom, so the buffer
		// must be large enough to read the entire reply packet. See:
		// https://groups.google.com/g/golang-nuts/c/5dy2Q4nPs08/m/KmuSQAGEtG4J
		bytesRead, _, err := conn.ReadFrom(buffer)
		if err != nil {
			return 0, fmt.Errorf("reading from ICMP connection: %w", err)
		}
		packetBytes := buffer[:bytesRead]

		packetBytes = packetBytes[ipv6.HeaderLen:]

		inboundMessage, err := icmp.ParseMessage(icmpv6Protocol, packetBytes)
		if err != nil {
			return 0, fmt.Errorf("parsing message: %w", err)
		}

		switch typedBody := inboundMessage.Body.(type) {
		case *icmp.PacketTooBig:
			// https://datatracker.ietf.org/doc/html/rfc1885#section-3.2
			mtu = typedBody.MTU
			err = checkMTU(mtu, minIPv6MTU, physicalLinkMTU)
			if err != nil {
				return 0, fmt.Errorf("checking MTU: %w", err)
			}

			// Sanity checks
			const truncatedBody = true
			err = checkEchoReply(icmpv6Protocol, typedBody.Data, outboundMessage, truncatedBody)
			if err != nil {
				return 0, fmt.Errorf("checking invoking message: %w", err)
			}
			return typedBody.MTU, nil
		case *icmp.DstUnreach:
			// https://datatracker.ietf.org/doc/html/rfc1885#section-3.1
			idMatch, err := checkInvokingReplyIDMatch(icmpv6Protocol, packetBytes, outboundMessage)
			if err != nil {
				return 0, fmt.Errorf("checking invoking message id: %w", err)
			} else if idMatch {
				return 0, fmt.Errorf("%w", ErrICMPDestinationUnreachable)
			}
			logger.Debug("discarding received ICMP destination unreachable reply with an unknown id")
			continue
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
