//go:build linux || darwin

package tcp

import (
	"syscall"
	"time"
)

// fileDescriptor is a platform-independent type for socket file descriptors.
type fileDescriptor int

func sendTo(fd fileDescriptor, p []byte, flags int, to syscall.Sockaddr) (err error) {
	return syscall.Sendto(int(fd), p, flags, to)
}

func setSocketTimeout(fd fileDescriptor, timeout time.Duration) (err error) {
	timeval := syscall.NsecToTimeval(timeout.Nanoseconds())
	return syscall.SetsockoptTimeval(int(fd), syscall.SOL_SOCKET, syscall.SO_RCVTIMEO, &timeval)
}

func recvFrom(fd fileDescriptor, p []byte, flags int) (n int, from syscall.Sockaddr, err error) {
	return syscall.Recvfrom(int(fd), p, flags)
}

func setNonBlock(fd int) error {
	return syscall.SetNonblock(fd, true)
}
