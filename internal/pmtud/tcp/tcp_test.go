package tcp

import (
	"context"
	"net/netip"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_runTest(t *testing.T) {
	t.Parallel()

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
			mtu:     1430,
			success: true,
		},
		"remote_not_listening": {
			dst: func(_ *testing.T) netip.AddrPort {
				return netip.AddrPortFrom(netip.AddrFrom4([4]byte{1, 1, 1, 1}), 12345)
			},
			mtu: 1300,
		},
		"1.1.1.1:443_mtu1300": {
			dst: func(_ *testing.T) netip.AddrPort {
				return netip.AddrPortFrom(netip.AddrFrom4([4]byte{1, 1, 1, 1}), 443)
			},
			mtu:     1300,
			success: true,
		},
		"1.1.1.1:80_mtu1400": {
			dst: func(_ *testing.T) netip.AddrPort {
				return netip.AddrPortFrom(netip.AddrFrom4([4]byte{1, 1, 1, 1}), 80)
			},
			mtu:     1400,
			success: true,
		},
		"1.1.1.1:80_mtu1480": {
			dst: func(_ *testing.T) netip.AddrPort {
				return netip.AddrPortFrom(netip.AddrFrom4([4]byte{1, 1, 1, 1}), 80)
			},
			mtu:     1480,
			success: true,
		},
		"8.8.8.8:443_mtu1300": {
			dst: func(_ *testing.T) netip.AddrPort {
				return netip.AddrPortFrom(netip.AddrFrom4([4]byte{8, 8, 8, 8}), 443)
			},
			mtu:     1300,
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
