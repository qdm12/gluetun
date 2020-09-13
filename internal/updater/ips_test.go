package updater

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_uniqueSortedIPs(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		inputIPs  []net.IP
		outputIPs []net.IP
	}{
		"nil": {
			inputIPs:  nil,
			outputIPs: []net.IP{},
		},
		"empty": {
			inputIPs:  []net.IP{},
			outputIPs: []net.IP{},
		},
		"single IPv4": {
			inputIPs:  []net.IP{{1, 1, 1, 1}},
			outputIPs: []net.IP{{1, 1, 1, 1}},
		},
		"two IPv4s": {
			inputIPs:  []net.IP{{1, 1, 2, 1}, {1, 1, 1, 1}},
			outputIPs: []net.IP{{1, 1, 1, 1}, {1, 1, 2, 1}},
		},
	}
	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			outputIPs := uniqueSortedIPs(testCase.inputIPs)
			assert.Equal(t, testCase.outputIPs, outputIPs)
		})
	}
}
