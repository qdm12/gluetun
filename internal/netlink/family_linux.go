package netlink

import "golang.org/x/sys/unix"

const (
	FamilyAll = unix.AF_UNSPEC
	FamilyV4  = unix.AF_INET
	FamilyV6  = unix.AF_INET6
)
