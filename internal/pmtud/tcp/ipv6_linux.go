package tcp

import "syscall"

func setIPv6HeaderIncluded(fd int) error {
	const ipv6HdrIncluded = 36 // IPV6_HDRINCL
	return syscall.SetsockoptInt(fd, syscall.IPPROTO_IPV6, ipv6HdrIncluded, 1)
}
