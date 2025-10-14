package pmtud

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

var (
	ErrICMPNotPermitted                            = errors.New("ICMP not permitted")
	ErrICMPDestinationUnreachable                  = errors.New("ICMP destination unreachable")
	ErrICMPCommunicationAdministrativelyProhibited = errors.New("communication administratively prohibited")
	ErrICMPBodyUnsupported                         = errors.New("ICMP body type is not supported")
)

func wrapConnErr(err error, timedCtx context.Context, pingTimeout time.Duration) error { //nolint:revive
	switch {
	case strings.HasSuffix(err.Error(), "sendto: operation not permitted"):
		err = fmt.Errorf("%w", ErrICMPNotPermitted)
	case errors.Is(timedCtx.Err(), context.DeadlineExceeded):
		err = fmt.Errorf("%w (timed out after %s)", net.ErrClosed, pingTimeout)
	case timedCtx.Err() != nil:
		err = timedCtx.Err()
	}
	return err
}
