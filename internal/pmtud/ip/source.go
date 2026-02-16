package ip

import (
	"fmt"
	"net/netip"
	"syscall"

	"github.com/jsimonetti/rtnetlink"
)

// SrcAddr determines the appropriate source IP address to use when sending a packet to the
// specified destination. It also reserves an ephemeral source port for the specified protocol
// to ensure that the port is not used by other processes. The cleanup function returned should
// be called to release the reserved port when done.
func SrcAddr(dst netip.AddrPort, proto int) (src netip.AddrPort, cleanup func(), err error) {
	srcAddr, err := srcIP(dst.Addr())
	if err != nil {
		return netip.AddrPort{}, nil, fmt.Errorf("finding source IP: %w", err)
	}

	srcPort, cleanup, err := srcPort(srcAddr, proto)
	if err != nil {
		return netip.AddrPort{}, nil, fmt.Errorf("reserving source port: %w", err)
	}

	return netip.AddrPortFrom(srcAddr, srcPort), cleanup, nil
}

var errNoRoute = fmt.Errorf("no route to destination")

func srcIP(dst netip.Addr) (netip.Addr, error) {
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

// srcPort reserves an ephemeral source port by opening a socket for the
// protocol specified and binds it to the provided source address.
// It doesn't actually listen on the port.
// The cleanup function returned should be called to release the port when done.
func srcPort(srcAddr netip.Addr, proto int) (srcPort uint16, cleanup func(), err error) {
	family := syscall.AF_INET
	if srcAddr.Is6() {
		family = syscall.AF_INET6
	}

	fd, err := syscall.Socket(family, syscall.SOCK_STREAM, proto)
	if err != nil {
		return 0, nil, fmt.Errorf("creating reservation socket: %w", err)
	}
	cleanup = func() {
		_ = syscall.Close(fd)
	}

	// Bind to port 0 to get an ephemeral port
	const port = 0
	var bindAddr syscall.Sockaddr
	if srcAddr.Is4() {
		bindAddr = &syscall.SockaddrInet4{
			Port: port,
			Addr: srcAddr.As4(),
		}
	} else {
		bindAddr = &syscall.SockaddrInet6{
			Port: port,
			Addr: srcAddr.As16(),
		}
	}

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
