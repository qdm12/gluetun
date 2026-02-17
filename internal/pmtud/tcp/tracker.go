package tcp

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/qdm12/gluetun/internal/pmtud/constants"
)

type tracker struct {
	fd              fileDescriptor
	ipv4            bool
	mutex           sync.RWMutex
	portsToDispatch map[uint32]dispatch
}

type dispatch struct {
	replyCh chan<- []byte
	abort   <-chan struct{}
}

func newTracker(fd fileDescriptor, ipv4 bool) *tracker {
	return &tracker{
		fd:              fd,
		ipv4:            ipv4,
		portsToDispatch: make(map[uint32]dispatch),
	}
}

func (t *tracker) constructKey(localPort, remotePort uint16) uint32 {
	buf := make([]byte, 4) //nolint:mnd
	binary.BigEndian.PutUint16(buf[0:2], localPort)
	binary.BigEndian.PutUint16(buf[2:4], remotePort)
	return binary.BigEndian.Uint32(buf)
}

func (t *tracker) register(localPort, remotePort uint16,
	ch chan<- []byte, abort <-chan struct{},
) {
	key := t.constructKey(localPort, remotePort)
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.portsToDispatch[key] = dispatch{
		replyCh: ch,
		abort:   abort,
	}
}

func (t *tracker) unregister(localPort, remotePort uint16) {
	key := t.constructKey(localPort, remotePort)
	t.mutex.Lock()
	defer t.mutex.Unlock()
	delete(t.portsToDispatch, key)
}

// listen listens for incoming TCP packets and dispatches them to the
// correct channel based on the source and destination port.
// If the context has a deadline associated, this one is used on the socket.
// Note it returns a nil error on context cancellation.
func (t *tracker) listen(ctx context.Context) error {
	deadline, hasDeadline := ctx.Deadline()
	for ctx.Err() == nil {
		if hasDeadline {
			remaining := time.Until(deadline)
			if remaining <= 0 {
				return nil
			}
			err := setSocketTimeout(t.fd, remaining)
			if err != nil {
				return fmt.Errorf("setting socket receive timeout: %w", err)
			}
		}

		reply := make([]byte, constants.MaxEthernetFrameSize)
		n, _, err := recvFrom(t.fd, reply, 0)
		if err != nil {
			switch {
			case errors.Is(err, constants.EAGAIN), errors.Is(err, constants.EWOULDBLOCK):
				pollSleep(ctx)
				continue
			case ctx.Err() != nil:
				// context canceled, stop listening so exit cleanly with no error
				return nil //nolint:nilerr
			default:
				return fmt.Errorf("receiving on socket: %w", err)
			}
		}
		reply = reply[:n]

		if t.ipv4 {
			var ok bool
			reply, ok = stripIPv4Header(reply)
			if !ok {
				continue // not an IPv4 packet
			}
		}

		const minTCPHeaderLength = 20
		if len(reply) < minTCPHeaderLength {
			continue
		}

		srcPort := binary.BigEndian.Uint16(reply[0:2])
		dstPort := binary.BigEndian.Uint16(reply[2:4])
		key := t.constructKey(dstPort, srcPort)
		t.mutex.RLock()
		dispatch, exists := t.portsToDispatch[key]
		t.mutex.RUnlock()
		if !exists {
			continue
		}
		select {
		case dispatch.replyCh <- reply:
		case <-dispatch.abort:
		}
	}
	return nil
}

func pollSleep(ctx context.Context) {
	const sleepBetweenPolls = 10 * time.Millisecond
	timer := time.NewTimer(sleepBetweenPolls)
	select {
	case <-ctx.Done():
		timer.Stop()
	case <-timer.C:
	}
}
