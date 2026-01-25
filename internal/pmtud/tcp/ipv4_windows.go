package tcp

import (
	"syscall"

	"golang.org/x/sys/windows"
)

func setIPv4HeaderIncluded(handle syscall.Handle) error {
	const ipHdrIncluded = windows.IP_HDRINCL
	return syscall.SetsockoptInt(handle, syscall.IPPROTO_IP, ipHdrIncluded, 1)
}
