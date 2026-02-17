package tcp

import (
	"net/netip"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

type fileDescriptor windows.Handle

func socket(domain int, typ int, proto int) (fd windows.Handle, err error) {
	return windows.Socket(domain, typ, proto)
}

func closeSocket(fd windows.Handle) error {
	return windows.Close(fd)
}

func sendTo(fd fileDescriptor, p []byte, flags int, to windows.Sockaddr) (err error) {
	return windows.Sendto(windows.Handle(fd), p, flags, to)
}

func setSocketTimeout(fd fileDescriptor, timeout time.Duration) (err error) {
	timeval := int(timeout.Milliseconds())
	return windows.SetsockoptInt(windows.Handle(fd), windows.SOL_SOCKET, windows.SO_RCVTIMEO, timeval)
}

func recvFrom(fd fileDescriptor, p []byte, flags int) (n int, from windows.Sockaddr, err error) {
	return windows.Recvfrom(windows.Handle(fd), p, flags)
}

func setMark(fd windows.Handle, _ int) error {
	panic("not implemented")
}

func setMTUDiscovery(fd windows.Handle) error {
	panic("not implemented")
}

func setNonBlock(fd windows.Handle) error {
	// Windows: Use ioctlsocket with FIONBIO
	var arg uint32 = 1 // 1 to enable non-blocking mode
	var bytesReturned uint32
	const FIONBIO = 0x8004667e
	return windows.WSAIoctl(fd, FIONBIO, (*byte)(unsafe.Pointer(&arg)),
		uint32(unsafe.Sizeof(arg)), nil, 0, &bytesReturned, nil, 0)
}

func makeSockAddr(addr netip.AddrPort) windows.Sockaddr {
	if addr.Addr().Is4() {
		return &windows.SockaddrInet4{
			Port: int(addr.Port()),
			Addr: addr.Addr().As4(),
		}
	}
	return &windows.SockaddrInet6{
		Port: int(addr.Port()),
		Addr: addr.Addr().As16(),
	}
}
