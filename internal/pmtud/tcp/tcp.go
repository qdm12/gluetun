package tcp

import (
	"context"
	"errors"
	"fmt"
	"net/netip"

	"github.com/qdm12/gluetun/internal/pmtud/constants"
	"github.com/qdm12/gluetun/internal/pmtud/ip"
)

func startRawSocket(family, excludeMark int) (fd fileDescriptor, stop func(), err error) {
	fdPlatform, err := socket(family, constants.SOCK_RAW, constants.IPPROTO_TCP)
	if err != nil {
		return 0, nil, fmt.Errorf("creating raw socket: %w", err)
	}

	err = setMark(fdPlatform, excludeMark)
	if err != nil {
		_ = closeSocket(fdPlatform)
		return 0, nil, fmt.Errorf("setting mark option on raw socket: %w", err)
	}

	if family == constants.AF_INET {
		err = ip.SetIPv4HeaderIncluded(fdPlatform)
	} else {
		err = ip.SetIPv6HeaderIncluded(fdPlatform)
	}
	if err != nil {
		_ = closeSocket(fdPlatform)
		return 0, nil, fmt.Errorf("setting header option on raw socket: %w", err)
	}

	// Allow sending packets larger than cached PMTU (for PMTUD probing)
	err = setMTUDiscovery(fdPlatform)
	if err != nil {
		_ = closeSocket(fdPlatform)
		return 0, nil, fmt.Errorf("setting IP_MTU_DISCOVER: %w", err)
	}

	// use polling because some Linux systems do not cancel
	// blocking syscalls such as recvfrom when the socket is closed,
	// which would cause things to hang indefinitely.
	err = setNonBlock(fdPlatform)
	if err != nil {
		_ = closeSocket(fdPlatform)
		return 0, nil, fmt.Errorf("setting non-blocking mode: %w", err)
	}

	stop = func() {
		_ = closeSocket(fdPlatform)
	}
	return fileDescriptor(fdPlatform), stop, nil
}

var (
	errTCPPacketNotSynAck        = errors.New("TCP packet is not a SYN-ACK")
	errTCPSynAckAckMismatch      = errors.New("TCP SYN-ACK ACK number does not match expected value")
	errFinalPacketTypeUnexpected = errors.New("final TCP packet type is unexpected")
	errTCPPacketLost             = errors.New("TCP packet was lost")
)

// Craft and send a raw TCP packet to test the MTU.
// It expects either an RST reply (if no server is listening)
// or a SYN-ACK/ACK reply (if a server is listening).
func runTest(ctx context.Context, dst netip.AddrPort, mtu uint32,
	excludeMark int, fd fileDescriptor, tracker *tracker,
	firewall Firewall, logger Logger,
) error {
	const proto = constants.IPPROTO_TCP
	src, cleanup, err := ip.SrcAddr(dst, proto)
	if err != nil {
		return fmt.Errorf("getting source address: %w", err)
	}
	defer cleanup()

	revert, err := firewall.TempDropOutputTCPRST(ctx, src, dst, excludeMark)
	if err != nil {
		return fmt.Errorf("temporarily dropping outgoing TCP RST packets: %w", err)
	}
	defer func() {
		// we don't want to skip reverting the firewall changes
		// even if the context is already expired, so we use a
		// background context here.
		err := revert(context.Background())
		if err != nil {
			logger.Warnf("reverting firewall changes: %s", err)
		}
	}()

	ch := make(chan []byte)
	abort := make(chan struct{})
	defer close(abort)
	tracker.register(src.Port(), dst.Port(), ch, abort)
	defer tracker.unregister(src.Port(), dst.Port())

	dstSockAddr := makeSockAddr(dst)

	synPacket, synSeq := createSYNPacket(src, dst, 0)
	const sendToFlags = 0
	err = sendTo(fd, synPacket, sendToFlags, dstSockAddr)
	if err != nil {
		return fmt.Errorf("sending SYN packet: %w", err)
	}

	var reply []byte
	select {
	case <-ctx.Done():
		_ = sendRST(fd, src, dst, synSeq+1)
		return ctx.Err()
	case reply = <-ch:
	}

	firstReplyHeader, err := parseTCPHeader(reply)
	switch {
	case err != nil:
		return fmt.Errorf("parsing first reply TCP header: %w", err)
	case firstReplyHeader.typ == packetTypeRST,
		firstReplyHeader.typ == packetTypeRSTACK:
		// server actively closed the connection, try sending a SYN with data
		return handleRSTReply(ctx, fd, ch, src, dst, mtu)
	case firstReplyHeader.typ != packetTypeSYNACK:
		return fmt.Errorf("%w: unexpected packet type %s", errTCPPacketNotSynAck, firstReplyHeader.typ)
	case firstReplyHeader.ack != synSeq+1:
		return fmt.Errorf("%w: expected %d, got %d", errTCPSynAckAckMismatch, synSeq+1, firstReplyHeader.ack)
	}

	if firstReplyHeader.options.mss != 0 {
		// If the server sent an MSS option, make sure our test packet is not larger than that MSS.
		tcpDataLength := getPayloadLength(mtu, dst) - constants.BaseTCPHeaderLength
		if tcpDataLength > uint32(firstReplyHeader.options.mss) {
			diff := tcpDataLength - uint32(firstReplyHeader.options.mss)
			minMTU := constants.MinIPv4MTU
			if dst.Addr().Is6() {
				minMTU = constants.MinIPv6MTU
			}
			diff = min(diff, mtu-minMTU)
			mtu -= diff
		}
	}

	// Send an ACK packet to finish the 3-way handshake, together with the
	// data to test the MTU, using TCP fast-open.
	ackPacket := createACKPacket(src, dst, firstReplyHeader.ack, firstReplyHeader.seq+1, mtu)
	err = sendTo(fd, ackPacket, sendToFlags, dstSockAddr)
	if err != nil {
		return fmt.Errorf("sending ACK packet: %w", err)
	}

	select {
	case <-ctx.Done():
		_ = sendRST(fd, src, dst, firstReplyHeader.ack)
		return ctx.Err()
	case reply = <-ch:
	}

	finalPacketHeader, err := parseTCPHeader(reply)
	if err != nil {
		return fmt.Errorf("parsing second reply TCP header: %w", err)
	}

	switch finalPacketHeader.typ { //nolint:exhaustive
	case packetTypeRST:
		return nil
	case packetTypeACK:
		err = sendRST(fd, src, dst, finalPacketHeader.ack)
		if err != nil {
			return fmt.Errorf("sending RST packet: %w", err)
		}
		return nil
	case packetTypeSYNACK: // server never received our MTU-test ACK packet
		return fmt.Errorf("%w: server responded with second SYN-ACK packet", errTCPPacketLost)
	default:
		_ = sendRST(fd, src, dst, finalPacketHeader.ack)
		return fmt.Errorf("%w: %s", errFinalPacketTypeUnexpected, finalPacketHeader.typ)
	}
}

var errTCPPacketNotRST = errors.New("TCP packet is not an RST")

func handleRSTReply(ctx context.Context, fd fileDescriptor, ch <-chan []byte,
	src, dst netip.AddrPort, mtu uint32,
) error {
	packet, synSeq := createSYNPacket(src, dst, mtu)
	const sendToFlags = 0
	err := sendTo(fd, packet, sendToFlags, makeSockAddr(dst))
	if err != nil {
		return fmt.Errorf("sending SYN MTU-test packet: %w", err)
	}

	var reply []byte
	select {
	case <-ctx.Done():
		_ = sendRST(fd, src, dst, synSeq+1)
		return ctx.Err() // timeout: the MTU test SYN packet was too big
	case reply = <-ch:
	}

	replyPacketHeader, err := parseTCPHeader(reply)
	if err != nil {
		return fmt.Errorf("parsing reply TCP header: %w", err)
	} else if replyPacketHeader.typ != packetTypeRST &&
		replyPacketHeader.typ != packetTypeRSTACK {
		return fmt.Errorf("%w: %s", errTCPPacketNotRST, replyPacketHeader.typ)
	}
	return nil
}

func sendRST(fd fileDescriptor, src, dst netip.AddrPort,
	previousACK uint32,
) error {
	seq := previousACK
	const ack = 0
	rstPacket := createRSTPacket(src, dst, seq, ack)
	const sendToFlags = 0
	return sendTo(fd, rstPacket, sendToFlags, makeSockAddr(dst))
}
