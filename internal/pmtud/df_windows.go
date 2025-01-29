package pmtud

import (
	"syscall"
)

func setDontFragment(fd uintptr) (err error) {
	// https://docs.microsoft.com/en-us/troubleshoot/windows/win32/header-library-requirement-socket-ipproto-ip
	// #define IP_DONTFRAGMENT        14     /* don't fragment IP datagrams */
	return syscall.SetsockoptInt(syscall.Handle(fd), syscall.IPPROTO_IP, 14, 1)
}
