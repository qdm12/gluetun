package tcp

import (
	"syscall"
	"time"
	"unsafe"

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

func setNonBlock(fd syscall.Handle) error {
	// Windows: Use ioctlsocket with FIONBIO
	var arg uint32 = 1 // 1 to enable non-blocking mode
	var bytesReturned uint32
	const FIONBIO = 0x8004667e
	return syscall.WSAIoctl(fd, FIONBIO, (*byte)(unsafe.Pointer(&arg)),
		uint32(unsafe.Sizeof(arg)), nil, 0, &bytesReturned, nil, 0)
}
