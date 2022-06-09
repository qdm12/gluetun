package updater

import (
	"net"
	"testing"

	"github.com/qdm12/gluetun/internal/constants/vpn"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/stretchr/testify/assert"
)

func Test_hostToServer_add(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		initialHTS  hostToServer
		host        string
		country     string
		city        string
		tcp         bool
		udp         bool
		expectedHTS hostToServer
	}{
		"empty host to server": {
			initialHTS: hostToServer{},
			host:       "host",
			country:    "country",
			city:       "city",
			tcp:        true,
			udp:        true,
			expectedHTS: hostToServer{
				"host": {
					VPN:      vpn.OpenVPN,
					Hostname: "host",
					Country:  "country",
					City:     "city",
					TCP:      true,
					UDP:      true,
				},
			},
		},
		"add server": {
			initialHTS: hostToServer{
				"existing host": {},
			},
			host:    "host",
			country: "country",
			city:    "city",
			tcp:     true,
			udp:     true,
			expectedHTS: hostToServer{
				"existing host": {},
				"host": models.Server{
					VPN:      vpn.OpenVPN,
					Hostname: "host",
					Country:  "country",
					City:     "city",
					TCP:      true,
					UDP:      true,
				},
			},
		},
		"extend existing server": {
			initialHTS: hostToServer{
				"host": models.Server{
					VPN:      vpn.OpenVPN,
					Hostname: "host",
					Country:  "country",
					City:     "city",
					TCP:      true,
				},
			},
			host:    "host",
			country: "country",
			city:    "city",
			tcp:     false,
			udp:     true,
			expectedHTS: hostToServer{
				"host": models.Server{
					VPN:      vpn.OpenVPN,
					Hostname: "host",
					Country:  "country",
					City:     "city",
					TCP:      true,
					UDP:      true,
				},
			},
		},
	}
	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			testCase.initialHTS.add(testCase.host, testCase.country, testCase.city, testCase.tcp, testCase.udp)
			assert.Equal(t, testCase.expectedHTS, testCase.initialHTS)
		})
	}
}

func Test_hostToServer_toHostsSlice(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		hts   hostToServer
		hosts []string
	}{
		"empty host to server": {
			hts:   hostToServer{},
			hosts: []string{},
		},
		"single host": {
			hts: hostToServer{
				"A": {},
			},
			hosts: []string{"A"},
		},
		"multiple hosts": {
			hts: hostToServer{
				"A": {},
				"B": {},
			},
			hosts: []string{"A", "B"},
		},
	}
	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			hosts := testCase.hts.toHostsSlice()
			assert.ElementsMatch(t, testCase.hosts, hosts)
		})
	}
}

func Test_hostToServer_adaptWithIPs(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		initialHTS  hostToServer
		hostToIPs   map[string][]net.IP
		expectedHTS hostToServer
	}{
		"create server": {
			initialHTS: hostToServer{},
			hostToIPs: map[string][]net.IP{
				"A": {{1, 2, 3, 4}},
			},
			expectedHTS: hostToServer{
				"A": models.Server{
					IPs: []net.IP{{1, 2, 3, 4}},
				},
			},
		},
		"add IPs to existing server": {
			initialHTS: hostToServer{
				"A": models.Server{
					Country: "country",
				},
			},
			hostToIPs: map[string][]net.IP{
				"A": {{1, 2, 3, 4}},
			},
			expectedHTS: hostToServer{
				"A": models.Server{
					Country: "country",
					IPs:     []net.IP{{1, 2, 3, 4}},
				},
			},
		},
		"remove server without IP": {
			initialHTS: hostToServer{
				"A": models.Server{
					Country: "country",
				},
			},
			hostToIPs:   map[string][]net.IP{},
			expectedHTS: hostToServer{},
		},
	}
	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			testCase.initialHTS.adaptWithIPs(testCase.hostToIPs)
			assert.Equal(t, testCase.expectedHTS, testCase.initialHTS)
		})
	}
}

func Test_hostToServer_toServersSlice(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		hts     hostToServer
		servers []models.Server
	}{
		"empty host to server": {
			hts:     hostToServer{},
			servers: []models.Server{},
		},
		"multiple servers": {
			hts: hostToServer{
				"A": {Country: "A"},
				"B": {Country: "B"},
			},
			servers: []models.Server{
				{Country: "A"},
				{Country: "B"},
			},
		},
	}
	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			servers := testCase.hts.toServersSlice()
			assert.ElementsMatch(t, testCase.servers, servers)
		})
	}
}
