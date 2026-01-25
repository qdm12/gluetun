package tcp

import (
	"context"
	"errors"
	"fmt"
	"net/netip"
	"syscall"

	"github.com/jsimonetti/rtnetlink"
	"github.com/qdm12/gluetun/internal/pmtud/constants"
)

func startRawSocket(family int) (fd fileDescriptor, stop func(), err error) {
	fdPlatform, err := syscall.Socket(family, syscall.SOCK_RAW, syscall.IPPROTO_TCP)
	if err != nil {
		return 0, nil, fmt.Errorf("creating raw socket: %w", err)
	}
	if family == syscall.AF_INET {
		err = setIPv4HeaderIncluded(fdPlatform)
	} else {
		err = setIPv6HeaderIncluded(fdPlatform)
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
	src, cleanup, err := getSrc(dst)
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

	ackPacket := createACKPacket(src, dst, synAckAck, synAckSeq+1, mtu)
	err = sendTo(fd, ackPacket, sendToFlags, dstSockAddr)
	if err != nil {
		return fmt.Errorf("sending ACK MTU-test packet: %w", err)
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
		return sendRST(fd, src, dst, ack)
	default:
		_ = sendRST(fd, src, dst, ack)
		return fmt.Errorf("%w: %s", errFinalPacketTypeUnexpected, packetType)
	}
}

func getSrc(dst netip.AddrPort) (src netip.AddrPort, cleanup func(), err error) {
	srcAddr, err := getSourceIP(dst.Addr())
	if err != nil {
		return netip.AddrPort{}, nil, fmt.Errorf("finding source IP: %w", err)
	}

	srcPort, cleanup, err := getSourcePort(srcAddr)
	if err != nil {
		return netip.AddrPort{}, nil, fmt.Errorf("reserving source port: %w", err)
	}

	return netip.AddrPortFrom(srcAddr, srcPort), cleanup, nil
}

var errNoRoute = fmt.Errorf("no route to destination")

func getSourceIP(dst netip.Addr) (netip.Addr, error) {
	conn, err := rtnetlink.Dial(nil)
	if err != nil {
		return netip.Addr{}, err
	}
	defer conn.Close()

	family := uint8(syscall.AF_INET)
	if dst.Is6() {
		family = syscall.AF_INET6
	}

	// Request route to destination
	requestMessage := &rtnetlink.RouteMessage{
		Family: family,
		Attributes: rtnetlink.RouteAttributes{
			Dst: dst.AsSlice(),
		},
	}
	messages, err := conn.Route.Get(requestMessage)
	if err != nil {
		return netip.Addr{}, fmt.Errorf("getting routes to %s: %w", dst, err)
	}

	for _, message := range messages {
		if message.Attributes.Src == nil {
			continue
		}
		ipv6 := message.Attributes.Src.To4() == nil
		if ipv6 {
			return netip.AddrFrom16([16]byte(message.Attributes.Src)), nil
		}
		return netip.AddrFrom4([4]byte(message.Attributes.Src)), nil
	}

	return netip.Addr{}, fmt.Errorf("%w: in %d route(s)", errNoRoute, len(messages))
}

// getSourcePort reserves an ephemeral source port by opening a TCP socket
// bound to the provided source address. It doesn't actually listen on the port.
// The cleanup function returned should be called to release the port when done.
func getSourcePort(srcAddr netip.Addr) (srcPort uint16, cleanup func(), err error) {
	family := syscall.AF_INET
	if srcAddr.Is6() {
		family = syscall.AF_INET6
	}

	fd, err := syscall.Socket(family, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		return 0, nil, fmt.Errorf("creating reservation socket: %w", err)
	}
	cleanup = func() {
		_ = syscall.Close(fd)
	}

	// Bind to port 0 to get an ephemeral port
	const port = 0
	addrPort := netip.AddrPortFrom(srcAddr, port)
	bindAddr := makeSockAddr(addrPort)

	err = syscall.Bind(fd, bindAddr)
	if err != nil {
		cleanup()
		return 0, nil, fmt.Errorf("binding reservation socket: %w", err)
	}

	sockAddr, err := syscall.Getsockname(fd)
	if err != nil {
		cleanup()
		return 0, nil, fmt.Errorf("getting bound socket name: %w", err)
	}

	switch typedSockAddr := sockAddr.(type) {
	case *syscall.SockaddrInet4:
		srcPort = uint16(typedSockAddr.Port) //nolint:gosec
	case *syscall.SockaddrInet6:
		srcPort = uint16(typedSockAddr.Port) //nolint:gosec
	default:
		panic(fmt.Sprintf("unexpected sockaddr type: %T", typedSockAddr))
	}
	return srcPort, cleanup, nil
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
	err := sendTo(fd, rstPacket, sendToFlags, makeSockAddr(dst))
	if err != nil {
		return fmt.Errorf("sending RST packet: %w", err)
	}
	return nil
}
