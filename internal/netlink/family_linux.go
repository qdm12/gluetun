package netlink

import "golang.org/x/sys/unix"

const (
	FamilyAll uint8 = unix.AF_UNSPEC
	FamilyV4  uint8 = unix.AF_INET
	FamilyV6  uint8 = unix.AF_INET6
)
