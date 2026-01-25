package tcp

import (
	"context"
	"errors"
	"fmt"
	"net/netip"
	"syscall"
	"time"

	"github.com/qdm12/gluetun/internal/pmtud/test"
)

var ErrMTUNotFound = errors.New("MTU not found")

func PathMTUDiscover(ctx context.Context, addr netip.AddrPort,
	minMTU, maxPossibleMTU uint32, logger Logger,
) (maxMTU uint32, err error) {
	return pmtudMultiSizes(ctx, addr, minMTU, maxPossibleMTU, logger)
}

type testUnit struct {
	mtu uint32
	ok  bool
}

func pmtudMultiSizes(ctx context.Context, addr netip.AddrPort,
	minMTU, maxPossibleMTU uint32, logger Logger,
) (maxMTU uint32, err error) {
	mtusToTest := test.MakeMTUsToTest(minMTU, maxPossibleMTU)
	if len(mtusToTest) == 1 { // only minMTU because minMTU == maxPossibleMTU
		return minMTU, nil
	}
	logger.Debugf("TCP testing the following MTUs: %v", mtusToTest)

	tests := make([]testUnit, len(mtusToTest))
	for i := range mtusToTest {
		tests[i] = testUnit{mtu: mtusToTest[i]}
	}

	family := syscall.AF_INET
	if addr.Addr().Is6() {
		family = syscall.AF_INET6
	}
	fd, stop, err := startRawSocket(family)
	if err != nil {
		return 0, fmt.Errorf("starting raw socket: %w", err)
	}
	defer stop()

	tracker := newTracker(fd, addr.Addr().Is4())

	const timeout = time.Second
	runCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	errCh := make(chan error)
	go func() {
		errCh <- tracker.listen(runCtx)
	}()

	doneCh := make(chan struct{})
	for i := range tests {
		go func(i int) {
			err := runTest(runCtx, fd, tracker, addr, tests[i].mtu)
			tests[i].ok = err == nil
			doneCh <- struct{}{}
		}(i)
	}

	for range tests {
		select {
		case <-doneCh:
		case err := <-errCh:
			if err == nil { // timeout
				break
			}
			return 0, fmt.Errorf("listening for TCP replies: %w", err)
		}
	}

	if tests[len(tests)-1].ok {
		return tests[len(tests)-1].mtu, nil
	}

	for i := len(tests) - 2; i >= 0; i-- { //nolint:mnd
		if tests[i].ok {
			stop()
			cancel()
			return pmtudMultiSizes(ctx, addr,
				tests[i].mtu, tests[i+1].mtu-1, logger)
		}
	}

	return 0, fmt.Errorf("%w: TCP queries to %s might be dropped "+
		"by the VPN server firewall", ErrMTUNotFound, addr)
}
