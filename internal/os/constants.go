package os

import (
	nativeos "os"
)

// Constants used for convenience so "os" does not have to be imported

//nolint:golint
const (
	O_CREATE = nativeos.O_CREATE
	O_TRUNC  = nativeos.O_TRUNC
	O_WRONLY = nativeos.O_WRONLY
	O_RDONLY = nativeos.O_RDONLY
	O_RDWR   = nativeos.O_RDWR
)
