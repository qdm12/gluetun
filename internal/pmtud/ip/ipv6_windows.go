package ip

import "syscall"

func SetIPv6HeaderIncluded(fd syscall.Handle) error {
	panic("windows does not allow an application to build IPv6 headers")
}
