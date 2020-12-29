package unix

import sysunix "golang.org/x/sys/unix"

// Constants used for convenience so "os" does not have to be imported

const (
	S_IFCHR = sysunix.S_IFCHR
)
