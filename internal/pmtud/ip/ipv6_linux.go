package ip

import "golang.org/x/sys/unix"

func SetIPv6HeaderIncluded(fd int) error {
	return unix.SetsockoptInt(fd, unix.IPPROTO_IPV6, unix.IPV6_HDRINCL, 1)
}
