//go:build linux || darwin

package ip

import "golang.org/x/sys/unix"

func SetIPv4HeaderIncluded(fd int) error {
	return unix.SetsockoptInt(fd, unix.IPPROTO_IP, unix.IP_HDRINCL, 1)
}
