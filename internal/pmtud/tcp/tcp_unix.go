//go:build linux || darwin

package tcp

import (
	"net/netip"
	"time"

	"golang.org/x/sys/unix"
)

// fileDescriptor is a platform-independent type for socket file descriptors.
type fileDescriptor int

func socket(domain int, typ int, proto int) (fd int, err error) {
	return unix.Socket(domain, typ, proto)
}

func closeSocket(fd int) error {
	return unix.Close(fd)
}

func sendTo(fd fileDescriptor, p []byte, flags int, to unix.Sockaddr) (err error) {
	return unix.Sendto(int(fd), p, flags, to)
}

func setSocketTimeout(fd fileDescriptor, timeout time.Duration) (err error) {
	timeval := unix.NsecToTimeval(timeout.Nanoseconds())
	return unix.SetsockoptTimeval(int(fd), unix.SOL_SOCKET, unix.SO_RCVTIMEO, &timeval)
}

func recvFrom(fd fileDescriptor, p []byte, flags int) (n int, from unix.Sockaddr, err error) {
	return unix.Recvfrom(int(fd), p, flags)
}

func setNonBlock(fd int) error {
	return unix.SetNonblock(fd, true)
}

func makeSockAddr(addr netip.AddrPort) unix.Sockaddr {
	if addr.Addr().Is4() {
		return &unix.SockaddrInet4{
			Port: int(addr.Port()),
			Addr: addr.Addr().As4(),
		}
	}
	return &unix.SockaddrInet6{
		Port: int(addr.Port()),
		Addr: addr.Addr().As16(),
	}
}
