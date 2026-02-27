//go:build linux || darwin

package constants

import "golang.org/x/sys/unix"

//nolint:revive
const (
	SOCK_RAW    = unix.SOCK_RAW
	SOCK_STREAM = unix.SOCK_STREAM
	AF_INET     = unix.AF_INET
	AF_INET6    = unix.AF_INET6
	IPPROTO_TCP = unix.IPPROTO_TCP
	EAGAIN      = unix.EAGAIN
	EWOULDBLOCK = unix.EWOULDBLOCK
)
