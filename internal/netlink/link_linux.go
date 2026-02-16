package netlink

import "golang.org/x/sys/unix"

const (
	DeviceTypeEthernet DeviceType = unix.ARPHRD_ETHER
	DeviceTypeLoopback DeviceType = unix.ARPHRD_LOOPBACK
	DeviceTypeNone     DeviceType = unix.ARPHRD_NONE

	iffUp = unix.IFF_UP
)
