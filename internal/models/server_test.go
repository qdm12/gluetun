package models

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Server_Equal(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		a     *Server
		b     Server
		equal bool
	}{
		"same IPs": {
			a: &Server{
				IPs: []net.IP{net.IPv4(1, 2, 3, 4)},
			},
			b: Server{
				IPs: []net.IP{net.IPv4(1, 2, 3, 4)},
			},
			equal: true,
		},
		"same IP strings": {
			a: &Server{
				IPs: []net.IP{net.IPv4(1, 2, 3, 4)},
			},
			b: Server{
				IPs: []net.IP{{1, 2, 3, 4}},
			},
			equal: true,
		},
		"different IPs": {
			a: &Server{
				IPs: []net.IP{{1, 2, 3, 4}, {2, 3, 4, 5}},
			},
			b: Server{
				IPs: []net.IP{{1, 2, 3, 4}, {1, 2, 3, 4}},
			},
		},
		"all fields equal": {
			a: &Server{
				VPN:         "vpn",
				Country:     "country",
				Region:      "region",
				City:        "city",
				ISP:         "isp",
				Owned:       true,
				Number:      1,
				ServerName:  "server_name",
				Hostname:    "hostname",
				TCP:         true,
				UDP:         true,
				OvpnX509:    "x509",
				RetroLoc:    "retroloc",
				MultiHop:    true,
				WgPubKey:    "wgpubkey",
				Free:        true,
				Stream:      true,
				PortForward: true,
				IPs:         []net.IP{net.IPv4(1, 2, 3, 4)},
				Keep:        true,
			},
			b: Server{
				VPN:         "vpn",
				Country:     "country",
				Region:      "region",
				City:        "city",
				ISP:         "isp",
				Owned:       true,
				Number:      1,
				ServerName:  "server_name",
				Hostname:    "hostname",
				TCP:         true,
				UDP:         true,
				OvpnX509:    "x509",
				RetroLoc:    "retroloc",
				MultiHop:    true,
				WgPubKey:    "wgpubkey",
				Free:        true,
				Stream:      true,
				PortForward: true,
				IPs:         []net.IP{net.IPv4(1, 2, 3, 4)},
				Keep:        true,
			},
			equal: true,
		},
		"different field": {
			a: &Server{
				VPN: "vpn",
			},
			b: Server{
				VPN: "other vpn",
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ipsOfANotNil := testCase.a.IPs != nil
			ipsOfBNotNil := testCase.b.IPs != nil

			equal := testCase.a.Equal(testCase.b)

			assert.Equal(t, testCase.equal, equal)

			// Ensure IPs field is not modified
			if ipsOfANotNil {
				assert.NotNil(t, testCase.a)
			}
			if ipsOfBNotNil {
				assert.NotNil(t, testCase.b)
			}
		})
	}
}
