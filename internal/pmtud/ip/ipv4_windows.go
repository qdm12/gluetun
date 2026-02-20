package ip

import (
	"golang.org/x/sys/windows"
)

func SetIPv4HeaderIncluded(handle windows.Handle) error {
	const ipHdrIncluded = windows.IP_HDRINCL
	return windows.SetsockoptInt(handle, windows.IPPROTO_IP, ipHdrIncluded, 1)
}
