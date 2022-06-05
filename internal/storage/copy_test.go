package storage

import (
	"net"
	"testing"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_copyServer(t *testing.T) {
	t.Parallel()

	server := models.Server{
		Country: "a",
		IPs:     []net.IP{{1, 2, 3, 4}},
	}

	serverCopy := copyServer(server)

	assert.Equal(t, server, serverCopy)
	// Check for mutation
	serverCopy.IPs[0][0] = 9
	assert.NotEqual(t, server, serverCopy)
}

func Test_copyIPs(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		toCopy []net.IP
		copied []net.IP
	}{
		"nil": {},
		"empty": {
			toCopy: []net.IP{},
			copied: []net.IP{},
		},
		"single IP": {
			toCopy: []net.IP{{1, 1, 1, 1}},
			copied: []net.IP{{1, 1, 1, 1}},
		},
		"two IPs": {
			toCopy: []net.IP{{1, 1, 1, 1}, {2, 2, 2, 2}},
			copied: []net.IP{{1, 1, 1, 1}, {2, 2, 2, 2}},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Reserver leading 9 for copy modifications below
			for _, ipToCopy := range testCase.toCopy {
				require.NotEqual(t, 9, ipToCopy[0])
			}

			copied := copyIPs(testCase.toCopy)

			assert.Equal(t, testCase.copied, copied)

			if len(copied) > 0 {
				original := testCase.toCopy[0][0]
				testCase.toCopy[0][0] = 9
				assert.NotEqual(t, 9, copied[0][0])
				testCase.toCopy[0][0] = original

				copied[0][0] = 9
				assert.NotEqual(t, 9, testCase.toCopy[0][0])
			}
		})
	}
}
