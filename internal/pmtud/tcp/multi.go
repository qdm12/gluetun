package tcp

import (
	"context"
	"errors"
	"fmt"
	"net/netip"
	"time"

	"github.com/qdm12/gluetun/internal/pmtud/constants"
	"github.com/qdm12/gluetun/internal/pmtud/test"
)

var ErrMTUNotFound = errors.New("MTU not found")

type testUnit struct {
	mtu uint32
	ok  bool
}

func PathMTUDiscover(ctx context.Context, dst netip.AddrPort,
	minMTU, maxPossibleMTU uint32, tryTimeout time.Duration,
	firewall Firewall, logger Logger,
) (mtu uint32, err error) {
	family := constants.AF_INET
	if dst.Addr().Is6() {
		family = constants.AF_INET6
	}
	const excludeMark = 4325
	fd, stop, err := startRawSocket(family, excludeMark)
	if err != nil {
		return 0, fmt.Errorf("starting raw socket: %w", err)
	}
	defer stop()

	tracker := newTracker(fd, dst.Addr().Is4())

	trackerCtx, trackerCancel := context.WithCancel(ctx)
	defer trackerCancel()
	trackerErrCh := make(chan error)
	go func() {
		trackerErrCh <- tracker.listen(trackerCtx)
	}()

	pmtudCtx, pmtudCancel := context.WithCancel(ctx)
	defer pmtudCancel()
	type result struct {
		mtu uint32
		err error
	}
	pmtudResultCh := make(chan result)
	go func() {
		mtu, err := pathMTUDiscover(pmtudCtx, fd, dst, minMTU, maxPossibleMTU,
			excludeMark, tryTimeout, tracker, firewall, logger)
		pmtudResultCh <- result{mtu: mtu, err: err}
	}()

	select {
	case err = <-trackerErrCh:
		pmtudCancel()
		<-pmtudResultCh
		return 0, fmt.Errorf("listening for TCP replies: %w", err)
	case res := <-pmtudResultCh:
		trackerCancel()
		<-trackerErrCh
		return res.mtu, res.err
	}
}

var errTimedOut = errors.New("timed out")

func pathMTUDiscover(ctx context.Context, fd fileDescriptor,
	dst netip.AddrPort, minMTU, maxPossibleMTU uint32, excludeMark int,
	tryTimeout time.Duration, tracker *tracker, firewall Firewall,
	logger Logger,
) (mtu uint32, err error) {
	mtusToTest := test.MakeMTUsToTest(minMTU, maxPossibleMTU)
	if len(mtusToTest) == 1 { // only minMTU because minMTU == maxPossibleMTU
		return minMTU, nil
	}
	logger.Debugf("TCP testing the following MTUs: %v", mtusToTest)

	tests := make([]testUnit, len(mtusToTest))
	for i := range mtusToTest {
		tests[i] = testUnit{mtu: mtusToTest[i]}
	}

	errCause := fmt.Errorf("%w: after %s", errTimedOut, tryTimeout)
	runCtx, runCancel := context.WithTimeoutCause(ctx, tryTimeout, errCause)
	defer runCancel()
	doneCh := make(chan struct{})
	for i := range tests {
		go func(i int) {
			err := runTest(runCtx, dst, tests[i].mtu, excludeMark,
				fd, tracker, firewall, logger)
			tests[i].ok = err == nil
			doneCh <- struct{}{}
		}(i)
	}

	i := 0
	for i < len(tests) {
		select {
		case <-runCtx.Done(): // timeout or parent context canceled
			err = context.Cause(runCtx)
			// collect remaining done signals
			for i < len(tests) {
				<-doneCh
				i++
			}
		case <-doneCh:
			i++
		}
	}

	if err != nil && !errors.Is(err, errTimedOut) {
		// context is canceled but did not timeout after tryTimeout
		return 0, fmt.Errorf("running MTU tests: %w", err)
	}

	if tests[len(tests)-1].ok {
		return tests[len(tests)-1].mtu, nil
	}

	for i := len(tests) - 2; i >= 0; i-- { //nolint:mnd
		if tests[i].ok {
			runCancel() // just to release resources although runCtx is no longer used
			return pathMTUDiscover(ctx, fd, dst,
				tests[i].mtu, tests[i+1].mtu-1, excludeMark,
				tryTimeout, tracker, firewall, logger)
		}
	}

	return 0, fmt.Errorf("%w: your connection might not be working at all", ErrMTUNotFound)
}
