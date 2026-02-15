package ip

import (
	"syscall"

	"golang.org/x/sys/windows"
)

func SetIPv4HeaderIncluded(handle syscall.Handle) error {
	const ipHdrIncluded = windows.IP_HDRINCL
	return syscall.SetsockoptInt(handle, syscall.IPPROTO_IP, ipHdrIncluded, 1)
}
