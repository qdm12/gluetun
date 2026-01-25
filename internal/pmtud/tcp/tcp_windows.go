package tcp

import (
	"syscall"
	"time"

	"golang.org/x/sys/windows"
)

type fileDescriptor syscall.Handle

func sendTo(fd fileDescriptor, p []byte, flags int, to syscall.Sockaddr) (err error) {
	return syscall.Sendto(syscall.Handle(fd), p, flags, to)
}

func setSocketTimeout(fd fileDescriptor, timeout time.Duration) (err error) {
	timeval := int(timeout.Milliseconds())
	return syscall.SetsockoptInt(syscall.Handle(fd), syscall.SOL_SOCKET, windows.SO_RCVTIMEO, timeval)
}

func recvFrom(fd fileDescriptor, p []byte, flags int) (n int, from syscall.Sockaddr, err error) {
	return syscall.Recvfrom(syscall.Handle(fd), p, flags)
}

func setMTUDiscovery(fd syscall.Handle) error {
	panic("not implemented")
}
