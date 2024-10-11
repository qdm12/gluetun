package perfectprivacy

import (
	"net/netip"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_internalIPToPorts(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		internalIP netip.Addr
		ports      []uint16
	}{
		"example_case": {
			internalIP: netip.AddrFrom4([4]byte{10, 0, 203, 88}),
			ports:      []uint16{12904, 22904, 32904},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ports := internalIPToPorts(testCase.internalIP)

			assert.Equal(t, testCase.ports, ports)
		})
	}
}
