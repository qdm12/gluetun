package resolver

import (
	"net/netip"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_uniqueIPsToSlice(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		inputIPs  map[string]struct{}
		outputIPs []netip.Addr
	}{
		"nil": {
			inputIPs:  nil,
			outputIPs: []netip.Addr{},
		},
		"empty": {
			inputIPs:  map[string]struct{}{},
			outputIPs: []netip.Addr{},
		},
		"single IPv4": {
			inputIPs:  map[string]struct{}{"1.1.1.1": {}},
			outputIPs: []netip.Addr{netip.AddrFrom4([4]byte{1, 1, 1, 1})},
		},
		"two IPv4s": {
			inputIPs:  map[string]struct{}{"1.1.1.1": {}, "1.1.2.1": {}},
			outputIPs: []netip.Addr{netip.AddrFrom4([4]byte{1, 1, 1, 1}), netip.AddrFrom4([4]byte{1, 1, 2, 1})},
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			outputIPs := uniqueIPsToSlice(testCase.inputIPs)
			assert.ElementsMatch(t, testCase.outputIPs, outputIPs)
		})
	}
}
