package storage

import (
	"net/netip"
	"testing"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/stretchr/testify/assert"
)

func Test_copyServer(t *testing.T) {
	t.Parallel()

	server := models.Server{
		Country: "a",
		IPs:     []netip.Addr{netip.AddrFrom4([4]byte{1, 2, 3, 4})},
	}

	serverCopy := copyServer(server)

	assert.Equal(t, server, serverCopy)
	// Check for mutation
	serverCopy.IPs[0] = netip.AddrFrom4([4]byte{9, 9, 9, 9})
	assert.NotEqual(t, server, serverCopy)
}

func Test_copyIPs(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		toCopy []netip.Addr
		copied []netip.Addr
	}{
		"nil": {},
		"empty": {
			toCopy: []netip.Addr{},
			copied: []netip.Addr{},
		},
		"single IP": {
			toCopy: []netip.Addr{netip.AddrFrom4([4]byte{1, 1, 1, 1})},
			copied: []netip.Addr{netip.AddrFrom4([4]byte{1, 1, 1, 1})},
		},
		"two IPs": {
			toCopy: []netip.Addr{netip.AddrFrom4([4]byte{1, 1, 1, 1}), netip.AddrFrom4([4]byte{2, 2, 2, 2})},
			copied: []netip.Addr{netip.AddrFrom4([4]byte{1, 1, 1, 1}), netip.AddrFrom4([4]byte{2, 2, 2, 2})},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			copied := copyIPs(testCase.toCopy)

			assert.Equal(t, testCase.copied, copied)

			if len(copied) > 0 {
				testCase.toCopy[0] = netip.AddrFrom4([4]byte{9, 9, 9, 9})
				assert.NotEqual(t, testCase.toCopy[0], testCase.copied[0])
			}
		})
	}
}
