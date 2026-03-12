package ip

import (
	"fmt"
	"net/netip"

	"golang.org/x/sys/windows"
)

func socket(domain int, typ int, proto int) (fd windows.Handle, err error) {
	return windows.Socket(domain, typ, proto)
}

func closeSocket(fd windows.Handle) error {
	return windows.Close(fd)
}

func bind(fd windows.Handle, addr windows.Sockaddr) error {
	return windows.Bind(fd, addr)
}

func makeSockAddr(ip netip.Addr, port uint16) windows.Sockaddr {
	if ip.Is4() {
		return &windows.SockaddrInet4{
			Port: int(port),
			Addr: ip.As4(),
		}
	}
	return &windows.SockaddrInet6{
		Port: int(port),
		Addr: ip.As16(),
	}
}

func extractPortFromFD(fd windows.Handle) (uint16, error) {
	sockAddr, err := windows.Getsockname(fd)
	if err != nil {
		return 0, fmt.Errorf("getting sockname: %w", err)
	}

	switch typedSockAddr := sockAddr.(type) {
	case *windows.SockaddrInet4:
		return uint16(typedSockAddr.Port), nil //nolint:gosec
	case *windows.SockaddrInet6:
		return uint16(typedSockAddr.Port), nil //nolint:gosec
	default:
		panic(fmt.Sprintf("unexpected sockaddr type: %T", typedSockAddr))
	}
}
