package pmtud

import (
	"errors"
)

var (
	ErrICMPDestinationUnreachable = errors.New("ICMP destination unreachable")
	ErrICMPBodyUnsupported        = errors.New("ICMP body type is not supported")
)
