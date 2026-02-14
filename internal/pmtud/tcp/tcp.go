package tcp

import (
	"context"
	"errors"
	"fmt"
	"net/netip"
	"syscall"

	"github.com/qdm12/gluetun/internal/pmtud/constants"
	"github.com/qdm12/gluetun/internal/pmtud/ip"
)

func startRawSocket(family int) (fd fileDescriptor, stop func(), err error) {
	fdPlatform, err := syscall.Socket(family, syscall.SOCK_RAW, syscall.IPPROTO_TCP)
	if err != nil {
		return 0, nil, fmt.Errorf("creating raw socket: %w", err)
	}

	if family == syscall.AF_INET {
		err = ip.SetIPv4HeaderIncluded(fdPlatform)
	} else {
		err = ip.SetIPv6HeaderIncluded(fdPlatform)
	}
	if err != nil {
		_ = syscall.Close(fdPlatform)
		return 0, nil, fmt.Errorf("setting header option on raw socket: %w", err)
	}

	// Allow sending packets larger than cached PMTU (for PMTUD probing)
	err = setMTUDiscovery(fdPlatform)
	if err != nil {
		_ = syscall.Close(fdPlatform)
		return 0, nil, fmt.Errorf("setting IP_MTU_DISCOVER: %w", err)
	}

	// use polling because some Linux systems do not cancel
	// blocking syscalls such as recvfrom when the socket is closed,
	// which would cause things to hang indefinitely.
	err = setNonBlock(fdPlatform)
	if err != nil {
		_ = syscall.Close(fdPlatform)
		return 0, nil, fmt.Errorf("setting non-blocking mode: %w", err)
	}

	stop = func() {
		_ = syscall.Close(fdPlatform)
	}
	return fileDescriptor(fdPlatform), stop, nil
}

var (
	errTCPPacketNotSynAck        = errors.New("TCP packet is not a SYN-ACK")
	errTCPSynAckAckMismatch      = errors.New("TCP SYN-ACK ACK number does not match expected value")
	errFinalPacketTypeUnexpected = errors.New("final TCP packet type is unexpected")
)

// Craft and send a raw TCP packet to test the MTU.
// It expects either an RST reply (if no server is listening)
// or a SYN-ACK/ACK reply (if a server is listening).
func runTest(ctx context.Context, fd fileDescriptor,
	tracker *tracker, dst netip.AddrPort, mtu uint32,
) error {
	const proto = syscall.IPPROTO_TCP
	src, cleanup, err := ip.SrcAddr(dst, proto)
	if err != nil {
		return fmt.Errorf("getting source address: %w", err)
	}
	defer cleanup()

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
		return ctx.Err()
	case reply = <-ch:
	}

	packetType, synAckSeq, synAckAck, err := parseTCPHeader(reply[:constants.BaseTCPHeaderLength])
	switch {
	case err != nil:
		return fmt.Errorf("parsing first reply TCP header: %w", err)
	case packetType == packetTypeRST:
		// server actively closed the connection, try sending a SYN with data
		return handleRSTReply(ctx, fd, ch, src, dst, mtu)
	case packetType != packetTypeSYNACK:
		return fmt.Errorf("%w: unexpected packet type %s", errTCPPacketNotSynAck, packetType)
	case synAckAck != synSeq+1:
		return fmt.Errorf("%w: expected %d, got %d", errTCPSynAckAckMismatch, synSeq+1, synAckAck)
	}

	// Send a no-data ACK packet to finish the 3-way handshake.
	const ackMTU = 0 // no data payload initially
	ackPacket := createACKPacket(src, dst, synAckAck, synAckSeq+1, ackMTU)
	err = sendTo(fd, ackPacket, sendToFlags, dstSockAddr)
	if err != nil {
		return fmt.Errorf("sending ACK-without-data packet: %w", err)
	}

	// Send a data ACK packet to test the MTU given.
	ackPacket = createACKPacket(src, dst, synAckAck, synAckSeq+1, mtu)
	err = sendTo(fd, ackPacket, sendToFlags, dstSockAddr)
	if err != nil {
		return fmt.Errorf("sending ACK-with-data packet: %w", err)
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case reply = <-ch:
	}

	packetType, _, ack, err := parseTCPHeader(reply[:constants.BaseTCPHeaderLength])
	if err != nil {
		return fmt.Errorf("parsing second reply TCP header: %w", err)
	}

	switch packetType { //nolint:exhaustive
	case packetTypeRST:
		return nil
	case packetTypeACK:
		err = sendRST(fd, src, dst, ack)
		if err != nil {
			return fmt.Errorf("sending RST packet: %w", err)
		}
		return nil
	default:
		_ = sendRST(fd, src, dst, ack)
		return fmt.Errorf("%w: %s", errFinalPacketTypeUnexpected, packetType)
	}
}

func makeSockAddr(addr netip.AddrPort) syscall.Sockaddr {
	if addr.Addr().Is4() {
		return &syscall.SockaddrInet4{
			Port: int(addr.Port()),
			Addr: addr.Addr().As4(),
		}
	}
	return &syscall.SockaddrInet6{
		Port: int(addr.Port()),
		Addr: addr.Addr().As16(),
	}
}

var errTCPPacketNotRST = errors.New("TCP packet is not an RST")

func handleRSTReply(ctx context.Context, fd fileDescriptor, ch <-chan []byte,
	src, dst netip.AddrPort, mtu uint32,
) error {
	packet, _ := createSYNPacket(src, dst, mtu)
	const sendToFlags = 0
	err := sendTo(fd, packet, sendToFlags, makeSockAddr(dst))
	if err != nil {
		return fmt.Errorf("sending SYN MTU-test packet: %w", err)
	}

	var reply []byte
	select {
	case <-ctx.Done():
		return ctx.Err() // timeout: the MTU test SYN packet was too big
	case reply = <-ch:
	}

	packetType, _, _, err := parseTCPHeader(reply[:constants.BaseTCPHeaderLength])
	if err != nil {
		return fmt.Errorf("parsing reply TCP header: %w", err)
	} else if packetType != packetTypeRST {
		return fmt.Errorf("%w: %s", errTCPPacketNotRST, packetType)
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
