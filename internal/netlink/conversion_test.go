package netlink

import (
	"net"
	"net/netip"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_netipPrefixToIPNet(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		prefix netip.Prefix
		ipNet  *net.IPNet
	}{
		"empty_prefix": {},
		"IPv4_prefix": {
			prefix: netip.PrefixFrom(netip.AddrFrom4([4]byte{1, 2, 3, 4}), 24),
			ipNet: &net.IPNet{
				IP:   net.IP{1, 2, 3, 4},
				Mask: net.IPv4Mask(255, 255, 255, 0),
			},
		},
		"IPv4-in-IPv6_prefix": {
			prefix: netip.PrefixFrom(netip.AddrFrom16(
				[16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff, 1, 2, 3, 4}),
				24),
			ipNet: &net.IPNet{
				IP:   net.IP{1, 2, 3, 4},
				Mask: net.IPv4Mask(255, 255, 255, 0),
			},
		},
		"IPv6_prefix": {
			prefix: netip.PrefixFrom(netip.IPv6Loopback(), 8),
			ipNet: &net.IPNet{
				IP:   net.IPv6loopback,
				Mask: net.IPMask{0xff, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ipNet := netipPrefixToIPNet(testCase.prefix)

			assert.Equal(t, testCase.ipNet, ipNet)
		})
	}
}

func Test_netIPNetToNetipPrefix(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		ipNet  *net.IPNet
		prefix netip.Prefix
	}{
		"empty ipnet": {},
		"custom sized IP in ipnet": {
			ipNet: &net.IPNet{
				IP: net.IP{1},
			},
		},
		"IPv4 ipnet": {
			ipNet: &net.IPNet{
				IP:   net.IP{1, 2, 3, 4},
				Mask: net.IPMask{255, 255, 255, 0},
			},
			prefix: netip.PrefixFrom(netip.AddrFrom4([4]byte{1, 2, 3, 4}), 24),
		},
		"IPv4-in-IPv6 ipnet": {
			ipNet: &net.IPNet{
				IP:   net.IPv4(1, 2, 3, 4),
				Mask: net.IPMask{255, 255, 255, 0},
			},
			prefix: netip.PrefixFrom(netip.AddrFrom4([4]byte{1, 2, 3, 4}), 24),
		},
		"IPv6 ipnet": {
			ipNet: &net.IPNet{
				IP:   net.IPv6loopback,
				Mask: net.IPMask{0xff},
			},
			prefix: netip.PrefixFrom(netip.IPv6Loopback(), 8),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			prefix := netIPNetToNetipPrefix(testCase.ipNet)

			assert.Equal(t, testCase.prefix, prefix)
		})
	}
}

func Test_netIPToNetipAddress(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		ip           net.IP
		address      netip.Addr
		panicMessage string
	}{
		"nil_ip": {},
		"ip_not_ipv4_or_ipv6": {
			ip: net.IP{1},
		},
		"IPv4": {
			ip:      net.IPv4(1, 2, 3, 4),
			address: netip.AddrFrom4([4]byte{1, 2, 3, 4}),
		},
		"IPv6": {
			ip:      net.IPv6zero,
			address: netip.AddrFrom16([16]byte{}),
		},
		"IPv4 prefixed with 0xffff": {
			ip:      net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff, 1, 2, 3, 4},
			address: netip.AddrFrom4([4]byte{1, 2, 3, 4}),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if testCase.panicMessage != "" {
				assert.PanicsWithValue(t, testCase.panicMessage, func() {
					netIPToNetipAddress(testCase.ip)
				})
				return
			}

			address := netIPToNetipAddress(testCase.ip)
			assert.Equal(t, testCase.address, address)
		})
	}
}
