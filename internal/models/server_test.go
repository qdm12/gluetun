package models

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PIAServer_String(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		server PIAServer
		s      string
	}{
		"no ips": {
			server: PIAServer{Region: "a b"},
			s:      `{Region: "a b", IPs: []net.IP{}}`,
		},
		"with ips": {
			server: PIAServer{Region: "a b", IPs: []net.IP{{1, 1, 1, 1}, {2, 2, 2, 2}}},
			s:      `{Region: "a b", IPs: []net.IP{{1, 1, 1, 1}, {2, 2, 2, 2}}}`,
		},
	}
	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			s := testCase.server.String()
			assert.Equal(t, testCase.s, s)
		})
	}
}

func Test_MullvadServer_String(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		server MullvadServer
		s      string
	}{
		"example": {
			server: MullvadServer{
				IPs:     []net.IP{{1, 1, 1, 1}},
				IPsV6:   []net.IP{{0x20, 0x1, 0xd, 0xb8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x0, 0x1}},
				Country: "That Country",
				City:    "That City",
				ISP:     "not spying on you",
				Owned:   true,
			},
			s: `{Country: "That Country", City: "That City", ISP: "not spying on you", Owned: true, IPs: []net.IP{{1, 1, 1, 1}}, IPsV6: []net.IP{{0x20, 0x1, 0xd, 0xb8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x0, 0x1}}}`,
		},
	}
	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			s := testCase.server.String()
			assert.Equal(t, testCase.s, s)
		})
	}
}

func Test_goStringifyIP(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		ip net.IP
		s  string
	}{
		"nil ip": {
			s: "net.IP{net.IP(nil)}",
		},
		"empty ip": {
			ip: net.IP{},
			s:  "net.IP{}",
		},
		"ipv4": {
			ip: net.IP{10, 16, 54, 25},
			s:  "net.IP{10, 16, 54, 25}",
		},
		"ipv6": {
			ip: net.IP{0x20, 0x1, 0xd, 0xb8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x0, 0x1},
			s:  "net.IP{0x20, 0x1, 0xd, 0xb8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x0, 0x1}",
		},
		"zeros ipv4": {
			ip: net.IP{0, 0, 0, 0},
			s:  "net.IP{}",
		},
		"zeros ipv46": {
			ip: net.ParseIP("::"),
			s:  "net.IP{}",
		},
	}
	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			s := goStringifyIP(testCase.ip)
			assert.Equal(t, testCase.s, s)
		})
	}
}

func Test_stringifyIPs(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		ips []net.IP
		s   string
	}{
		"nil ips": {
			s: "[]net.IP{}",
		},
		"empty ips": {
			ips: []net.IP{},
			s:   "[]net.IP{}",
		},
		"single ipv4": {
			ips: []net.IP{{10, 16, 54, 25}},
			s:   "[]net.IP{{10, 16, 54, 25}}",
		},
		"single ipv6": {
			ips: []net.IP{{0x20, 0x1, 0xd, 0xb8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x0, 0x1}},
			s:   "[]net.IP{{0x20, 0x1, 0xd, 0xb8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x0, 0x1}}",
		},
		"mix of ips": {
			ips: []net.IP{
				{10, 16, 54, 25},
				{0x20, 0x1, 0xd, 0xb8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x0, 0x1},
				{0, 0, 0, 0},
			},
			s: "[]net.IP{{10, 16, 54, 25}, {0x20, 0x1, 0xd, 0xb8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x0, 0x1}, {}}",
		},
	}
	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			s := goStringifyIPs(testCase.ips)
			assert.Equal(t, testCase.s, s)
		})
	}
}
