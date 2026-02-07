package tcp

import "syscall"

func setIPv6HeaderIncluded(fd syscall.Handle) error {
	panic("windows does not allow an application to build IPv6 headers")
}
