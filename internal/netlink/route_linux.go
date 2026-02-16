package netlink

import "golang.org/x/sys/unix"

const (
	RouteTypeUnicast = unix.RTN_UNICAST
	ScopeUniverse    = unix.RT_SCOPE_UNIVERSE
	ProtoStatic      = unix.RTPROT_STATIC

	rtTableCompat = unix.RT_TABLE_COMPAT
)
