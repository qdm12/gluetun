package routing

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_IPIsPrivate(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		ipString  string
		isPrivate bool
	}{
		"loopback 127.0.0.1": {
			ipString:  "127.0.0.1",
			isPrivate: true,
		},
		"loopback 127.0.0.10": {
			ipString:  "127.0.0.10",
			isPrivate: true,
		},
		"loopback ::1": {
			ipString:  "::1",
			isPrivate: true,
		},
		"private 10.0.0.1": {
			ipString:  "10.0.0.1",
			isPrivate: true,
		},
		"private 10.255.255.255": {
			ipString:  "10.255.255.255",
			isPrivate: true,
		},
		"private 172.16.0.1": {
			ipString:  "172.16.0.1",
			isPrivate: true,
		},
		"private 172.31.255.255": {
			ipString:  "172.31.255.255",
			isPrivate: true,
		},
		"private 192.168.0.0": {
			ipString:  "192.168.0.0",
			isPrivate: true,
		},
		"private 192.168.255.255": {
			ipString:  "192.168.255.255",
			isPrivate: true,
		},
		"private fc00::": {
			ipString:  "fc00::",
			isPrivate: true,
		},
		"private fc00::af": {
			ipString:  "fc00::af",
			isPrivate: true,
		},
		"local unicast 169.254.0.0": {
			ipString:  "169.254.0.0",
			isPrivate: true,
		},
		"local unicast 169.254.255.255": {
			ipString:  "169.254.255.255",
			isPrivate: true,
		},
		"local unicast fe80::": {
			ipString:  "fe80::",
			isPrivate: true,
		},
		"local unicast fe80::ae": {
			ipString:  "fe80::ae",
			isPrivate: true,
		},
		"public IPv4": {
			ipString: "11.5.6.7",
		},
		"public IPv6": {
			ipString: "af6d::",
		},
	}
	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ip := net.ParseIP(testCase.ipString)
			require.NotNil(t, ip)

			isPrivate := IPIsPrivate(ip)

			assert.Equal(t, testCase.isPrivate, isPrivate)
		})
	}
}

func Test_ensureNoIPv6WrappedIPv4(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		candidateIP net.IP
		resultIP    net.IP
	}{
		"nil": {},
		"ipv6": {
			candidateIP: net.IPv6loopback,
			resultIP:    net.IPv6loopback,
		},
		"ipv4": {
			candidateIP: net.IP{1, 2, 3, 4},
			resultIP:    net.IP{1, 2, 3, 4},
		},
		"ipv6_wrapped_ipv4": {
			candidateIP: net.IPv4(1, 2, 3, 4),
			resultIP:    net.IP{1, 2, 3, 4},
		},
	}
	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resultIP := ensureNoIPv6WrappedIPv4(testCase.candidateIP)
			assert.Equal(t, testCase.resultIP, resultIP)
		})
	}
}
