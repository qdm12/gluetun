package tcp

import (
	"context"
	"errors"
	"fmt"
	"net/netip"
	"syscall"
	"testing"
	"time"

	"github.com/qdm12/gluetun/internal/netlink"
	"github.com/qdm12/gluetun/internal/routing"
	"github.com/qdm12/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_runTest(t *testing.T) {
	t.Parallel()

	noopLogger := &noopLogger{}
	netlinker := netlink.New(noopLogger)
	loopbackMTU, err := findLoopbackMTU(netlinker)
	require.NoError(t, err, "finding loopback IPv4 MTU")
	defaultIPv4MTU, err := findDefaultIPv4RouteMTU(netlinker)
	require.NoError(t, err, "finding default IPv4 route MTU")
	const safetyMTUMargin = 0

	const timeout = 100 * time.Millisecond
	ctx, cancel := context.WithTimeout(t.Context(), timeout)
	t.Cleanup(cancel)

	const family = syscall.AF_INET
	fd, stop, err := startRawSocket(family)
	require.NoError(t, err)
	t.Cleanup(stop)

	const ipv4 = true
	tracker := newTracker(fd, ipv4)
	trackerCh := make(chan error)
	go func() {
		trackerCh <- tracker.listen(ctx)
	}()

	t.Cleanup(func() {
		cancel() // stop listening
		err = <-trackerCh
		require.NoError(t, err)
	})

	testCases := map[string]struct {
		dst     func(t *testing.T) netip.AddrPort
		mtu     uint32
		success bool
	}{
		"local_not_listening": {
			dst: func(t *testing.T) netip.AddrPort {
				t.Helper()
				port := reserveClosedPort(t)
				return netip.AddrPortFrom(netip.AddrFrom4([4]byte{127, 0, 0, 1}), port)
			},
			mtu:     loopbackMTU,
			success: true,
		},
		"remote_not_listening": {
			dst: func(_ *testing.T) netip.AddrPort {
				return netip.AddrPortFrom(netip.AddrFrom4([4]byte{1, 1, 1, 1}), 12345)
			},
			mtu: defaultIPv4MTU - safetyMTUMargin,
		},
		"1.1.1.1:443": {
			dst: func(_ *testing.T) netip.AddrPort {
				return netip.AddrPortFrom(netip.AddrFrom4([4]byte{1, 1, 1, 1}), 443)
			},
			mtu:     defaultIPv4MTU - safetyMTUMargin,
			success: true,
		},
		"1.1.1.1:80": {
			dst: func(_ *testing.T) netip.AddrPort {
				return netip.AddrPortFrom(netip.AddrFrom4([4]byte{1, 1, 1, 1}), 80)
			},
			mtu:     defaultIPv4MTU - safetyMTUMargin,
			success: true,
		},
		"8.8.8.8:443": {
			dst: func(_ *testing.T) netip.AddrPort {
				return netip.AddrPortFrom(netip.AddrFrom4([4]byte{8, 8, 8, 8}), 443)
			},
			mtu:     defaultIPv4MTU - safetyMTUMargin,
			success: true,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			dst := testCase.dst(t)
			err := runTest(ctx, fd, tracker, dst, testCase.mtu)
			if testCase.success {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

var errRouteNotFound = errors.New("route not found")

func findLoopbackMTU(netlinker *netlink.NetLink) (mtu uint32, err error) {
	routes, err := netlinker.RouteList(netlink.FamilyV4)
	if err != nil {
		return 0, fmt.Errorf("getting routes list: %w", err)
	}
	for _, route := range routes {
		if route.Dst.IsValid() && route.Dst.Addr().IsLoopback() {
			link, err := netlinker.LinkByIndex(route.LinkIndex)
			if err != nil {
				return 0, fmt.Errorf("getting link by index: %w", err)
			}
			// Quirk: make sure it is maximum 65535, and not i.e. 65536
			// or the IP header 16 bits will fail to fit that packet length value.
			const maxMTU = 65535
			return min(link.MTU, maxMTU), nil
		}
	}
	return 0, fmt.Errorf("%w: no loopback route found", errRouteNotFound)
}

func findDefaultIPv4RouteMTU(netlinker *netlink.NetLink) (mtu uint32, err error) {
	noopLogger := &noopLogger{}
	routing := routing.New(netlinker, noopLogger)
	defaultRoutes, err := routing.DefaultRoutes()
	if err != nil {
		return 0, fmt.Errorf("getting default routes: %w", err)
	}
	for _, route := range defaultRoutes {
		if route.Family != netlink.FamilyV4 {
			continue
		}
		link, err := netlinker.LinkByName(defaultRoutes[0].NetInterface)
		if err != nil {
			return 0, fmt.Errorf("getting link by name: %w", err)
		}
		return link.MTU, nil
	}
	return 0, fmt.Errorf("%w: no default route found", errRouteNotFound)
}

func reserveClosedPort(t *testing.T) (port uint16) {
	t.Helper()

	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	require.NoError(t, err)
	t.Cleanup(func() {
		err := syscall.Close(fd)
		assert.NoError(t, err)
	})

	addr := &syscall.SockaddrInet4{
		Port: 0,
		Addr: [4]byte{127, 0, 0, 1},
	}

	err = syscall.Bind(fd, addr)
	if err != nil {
		_ = syscall.Close(fd)
		t.Fatal(err)
	}

	sockAddr, err := syscall.Getsockname(fd)
	if err != nil {
		_ = syscall.Close(fd)
		t.Fatal(err)
	}

	sockAddr4, ok := sockAddr.(*syscall.SockaddrInet4)
	if !ok {
		_ = syscall.Close(fd)
		t.Fatal("not an IPv4 address")
	}

	return uint16(sockAddr4.Port) //nolint:gosec
}

type noopLogger struct{}

func (l *noopLogger) Patch(_ ...log.Option)     {}
func (l *noopLogger) Debug(_ string)            {}
func (l *noopLogger) Debugf(_ string, _ ...any) {}
func (l *noopLogger) Info(_ string)             {}
func (l *noopLogger) Warn(_ string)             {}
func (l *noopLogger) Error(_ string)            {}
