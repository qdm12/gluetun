package routing

import "golang.org/x/sys/unix"

const (
	tableMain  = unix.RT_TABLE_MAIN
	tableLocal = unix.RT_TABLE_LOCAL
)
