package icmp

import (
	"golang.org/x/sys/windows"
)

func setDontFragment(fd uintptr, ipv4 bool) (err error) {
	if ipv4 {
		// https://docs.microsoft.com/en-us/troubleshoot/windows/win32/header-library-requirement-socket-ipproto-ip
		// #define IP_DONTFRAGMENT        14     /* don't fragment IP datagrams */
		return windows.SetsockoptInt(windows.Handle(fd), windows.IPPROTO_IP, 14, 1)
	}
	return windows.SetsockoptInt(windows.Handle(fd), windows.IPPROTO_IPV6, 14, 1)
}
