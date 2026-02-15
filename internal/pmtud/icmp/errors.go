package icmp

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

var (
	ErrNotPermitted                            = errors.New("ICMP not permitted")
	ErrDestinationUnreachable                  = errors.New("ICMP destination unreachable")
	ErrCommunicationAdministrativelyProhibited = errors.New("communication administratively prohibited")
	ErrBodyUnsupported                         = errors.New("ICMP body type is not supported")
	ErrMTUNotFound                             = errors.New("MTU not found")
)

func wrapConnErr(err error, timedCtx context.Context, pingTimeout time.Duration) error { //nolint:revive
	switch {
	case strings.HasSuffix(err.Error(), "sendto: operation not permitted"):
		err = fmt.Errorf("%w", ErrNotPermitted)
	case errors.Is(timedCtx.Err(), context.DeadlineExceeded):
		err = fmt.Errorf("%w (timed out after %s)", net.ErrClosed, pingTimeout)
	case timedCtx.Err() != nil:
		err = timedCtx.Err()
	}
	return err
}
