package resolver

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_uniqueIPsToSlice(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		inputIPs  map[string]struct{}
		outputIPs []net.IP
	}{
		"nil": {
			inputIPs:  nil,
			outputIPs: []net.IP{},
		},
		"empty": {
			inputIPs:  map[string]struct{}{},
			outputIPs: []net.IP{},
		},
		"single IPv4": {
			inputIPs:  map[string]struct{}{"1.1.1.1": {}},
			outputIPs: []net.IP{{1, 1, 1, 1}},
		},
		"two IPv4s": {
			inputIPs:  map[string]struct{}{"1.1.1.1": {}, "1.1.2.1": {}},
			outputIPs: []net.IP{{1, 1, 1, 1}, {1, 1, 2, 1}},
		},
	}
	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			outputIPs := uniqueIPsToSlice(testCase.inputIPs)
			assert.ElementsMatch(t, testCase.outputIPs, outputIPs)
		})
	}
}
