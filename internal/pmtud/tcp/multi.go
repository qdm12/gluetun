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

func PathMTUDiscover(ctx context.Context, addrPort netip.AddrPort,
	minMTU, maxPossibleMTU uint32, firewall Firewall, logger Logger,
) (mtu uint32, err error) {
	const excludeMark = 4325
	revert, err := firewall.TempDropOutputTCPRST(ctx, addrPort, excludeMark)
	if err != nil {
		return 0, fmt.Errorf("temporarily dropping outgoing TCP RST packets: %w", err)
	}
	defer func() {
		err := revert(ctx)
		if err != nil {
			logger.Warnf("reverting firewall changes: %s", err)
		}
	}()

	return pathMTUDiscover(ctx, addrPort, minMTU, maxPossibleMTU, excludeMark, logger)
}

func pathMTUDiscover(ctx context.Context, addrPort netip.AddrPort,
	minMTU, maxPossibleMTU uint32, excludeMark int, logger Logger,
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

	family := constants.AF_INET
	if addrPort.Addr().Is6() {
		family = constants.AF_INET6
	}
	fd, stop, err := startRawSocket(family, excludeMark)
	if err != nil {
		return 0, fmt.Errorf("starting raw socket: %w", err)
	}
	defer stop()

	tracker := newTracker(fd, addrPort.Addr().Is4())

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
			err := runTest(runCtx, fd, tracker, src, dst, tests[i].mtu)
			tests[i].ok = err == nil
			doneCh <- struct{}{}
		}(i)
	}

	i := 0
	for i < len(tests) {
		select {
		case <-doneCh:
			i++
		case err := <-errCh:
			if err == nil { // timeout
				cancel()
				continue
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
			return pathMTUDiscover(ctx, addrPort,
				tests[i].mtu, tests[i+1].mtu-1, excludeMark, logger)
		}
	}

	return 0, fmt.Errorf("%w: your connection might not be working at all", ErrMTUNotFound)
}
