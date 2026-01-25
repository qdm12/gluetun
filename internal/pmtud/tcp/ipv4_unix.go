//go:build linux || darwin

package tcp

import "syscall"

func setIPv4HeaderIncluded(fd int) error {
	return syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_HDRINCL, 1)
}
