package updater

import (
	"net/netip"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_uniqueSortedIPs(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		inputIPs  []netip.Addr
		outputIPs []netip.Addr
	}{
		"nil": {
			inputIPs:  nil,
			outputIPs: []netip.Addr{},
		},
		"empty": {
			inputIPs:  []netip.Addr{},
			outputIPs: []netip.Addr{},
		},
		"single IPv4": {
			inputIPs:  []netip.Addr{netip.AddrFrom4([4]byte{1, 1, 1, 1})},
			outputIPs: []netip.Addr{netip.AddrFrom4([4]byte{1, 1, 1, 1})},
		},
		"two IPv4s": {
			inputIPs:  []netip.Addr{netip.AddrFrom4([4]byte{1, 1, 2, 1}), netip.AddrFrom4([4]byte{1, 1, 1, 1})},
			outputIPs: []netip.Addr{netip.AddrFrom4([4]byte{1, 1, 1, 1}), netip.AddrFrom4([4]byte{1, 1, 2, 1})},
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			outputIPs := uniqueSortedIPs(testCase.inputIPs)
			assert.Equal(t, testCase.outputIPs, outputIPs)
		})
	}
}
