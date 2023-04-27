package routing

import (
	"net"
	"net/netip"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_netIPToNetipAddress(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		ip           net.IP
		address      netip.Addr
		panicMessage string
	}{
		"nil ip": {
			panicMessage: "converting net.IP(nil) to netip.Addr failed",
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
		testCase := testCase
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
