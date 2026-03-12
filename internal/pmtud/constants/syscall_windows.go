package constants

import "golang.org/x/sys/windows"

const (
	SOCK_RAW    = windows.SOCK_RAW
	SOCK_STREAM = windows.SOCK_STREAM
	AF_INET     = windows.AF_INET
	AF_INET6    = windows.AF_INET6
	IPPROTO_TCP = windows.IPPROTO_TCP
	EAGAIN      = windows.WSAEWOULDBLOCK
	EWOULDBLOCK = windows.WSAEWOULDBLOCK
)
