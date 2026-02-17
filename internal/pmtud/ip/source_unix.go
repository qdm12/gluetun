//go:build linux || darwin

package ip

import (
	"fmt"
	"net/netip"

	"golang.org/x/sys/unix"
)

func socket(domain int, typ int, proto int) (fd int, err error) {
	return unix.Socket(domain, typ, proto)
}

func closeSocket(fd int) error {
	return unix.Close(fd)
}

func bind(fd int, addr unix.Sockaddr) error {
	return unix.Bind(fd, addr)
}

func makeSockAddr(ip netip.Addr, port uint16) unix.Sockaddr {
	if ip.Is4() {
		return &unix.SockaddrInet4{
			Port: int(port),
			Addr: ip.As4(),
		}
	}
	return &unix.SockaddrInet6{
		Port: int(port),
		Addr: ip.As16(),
	}
}

func extractPortFromFD(fd int) (uint16, error) {
	sockAddr, err := unix.Getsockname(fd)
	if err != nil {
		return 0, fmt.Errorf("getting sockname: %w", err)
	}

	switch typedSockAddr := sockAddr.(type) {
	case *unix.SockaddrInet4:
		return uint16(typedSockAddr.Port), nil //nolint:gosec
	case *unix.SockaddrInet6:
		return uint16(typedSockAddr.Port), nil //nolint:gosec
	default:
		panic(fmt.Sprintf("unexpected sockaddr type: %T", typedSockAddr))
	}
}
