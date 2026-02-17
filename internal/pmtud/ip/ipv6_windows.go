package ip

import "golang.org/x/sys/windows"

func SetIPv6HeaderIncluded(fd windows.Handle) error {
	panic("windows does not allow an application to build IPv6 headers")
}
