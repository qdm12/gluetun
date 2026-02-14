//go:build linux || darwin

package ip

import "syscall"

func SetIPv4HeaderIncluded(fd int) error {
	return syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_HDRINCL, 1)
}
