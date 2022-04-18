package ivpn

import (
	"errors"
	"math/rand"
	"testing"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/constants/vpn"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func boolPtr(b bool) *bool { return &b }

func Test_Ivpn_filterServers(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		servers   []models.Server
		selection settings.ServerSelection
		filtered  []models.Server
		err       error
	}{
		"no server available": {
			selection: settings.ServerSelection{}.WithDefaults(providers.Ivpn),
			err:       errors.New("no server found: for VPN openvpn; protocol udp"),
		},
		"no filter": {
			servers: []models.Server{
				{VPN: vpn.OpenVPN, Hostname: "a", UDP: true},
				{VPN: vpn.OpenVPN, Hostname: "b", UDP: true},
				{VPN: vpn.OpenVPN, Hostname: "c", UDP: true},
			},
			selection: settings.ServerSelection{}.WithDefaults(providers.Ivpn),
			filtered: []models.Server{
				{VPN: vpn.OpenVPN, Hostname: "a", UDP: true},
				{VPN: vpn.OpenVPN, Hostname: "b", UDP: true},
				{VPN: vpn.OpenVPN, Hostname: "c", UDP: true},
			},
		},
		"filter by country": {
			selection: settings.ServerSelection{
				Countries: []string{"b"},
			}.WithDefaults(providers.Ivpn),
			servers: []models.Server{
				{VPN: vpn.OpenVPN, Country: "a", UDP: true},
				{VPN: vpn.OpenVPN, Country: "b", UDP: true},
				{VPN: vpn.OpenVPN, Country: "c", UDP: true},
			},
			filtered: []models.Server{
				{VPN: vpn.OpenVPN, Country: "b", UDP: true},
			},
		},
		"filter by city": {
			selection: settings.ServerSelection{
				Cities: []string{"b"},
			}.WithDefaults(providers.Ivpn),
			servers: []models.Server{
				{VPN: vpn.OpenVPN, City: "a", UDP: true},
				{VPN: vpn.OpenVPN, City: "b", UDP: true},
				{VPN: vpn.OpenVPN, City: "c", UDP: true},
			},
			filtered: []models.Server{
				{VPN: vpn.OpenVPN, City: "b", UDP: true},
			},
		},
		"filter by ISP": {
			selection: settings.ServerSelection{
				ISPs: []string{"b"},
			}.WithDefaults(providers.Ivpn),
			servers: []models.Server{
				{VPN: vpn.OpenVPN, ISP: "a", UDP: true},
				{VPN: vpn.OpenVPN, ISP: "b", UDP: true},
				{VPN: vpn.OpenVPN, ISP: "c", UDP: true},
			},
			filtered: []models.Server{
				{VPN: vpn.OpenVPN, ISP: "b", UDP: true},
			},
		},
		"filter by hostname": {
			selection: settings.ServerSelection{
				Hostnames: []string{"b"},
			}.WithDefaults(providers.Ivpn),
			servers: []models.Server{
				{VPN: vpn.OpenVPN, Hostname: "a", UDP: true},
				{VPN: vpn.OpenVPN, Hostname: "b", UDP: true},
				{VPN: vpn.OpenVPN, Hostname: "c", UDP: true},
			},
			filtered: []models.Server{
				{VPN: vpn.OpenVPN, Hostname: "b", UDP: true},
			},
		},
		"filter by protocol": {
			selection: settings.ServerSelection{
				OpenVPN: settings.OpenVPNSelection{
					TCP: boolPtr(true),
				},
			}.WithDefaults(providers.Ivpn),
			servers: []models.Server{
				{VPN: vpn.OpenVPN, Hostname: "a", UDP: true},
				{VPN: vpn.OpenVPN, Hostname: "b", UDP: true, TCP: true},
				{VPN: vpn.OpenVPN, Hostname: "c", UDP: true},
			},
			filtered: []models.Server{
				{VPN: vpn.OpenVPN, Hostname: "b", UDP: true, TCP: true},
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			randSource := rand.NewSource(0)

			m := New(testCase.servers, randSource)

			servers, err := m.filterServers(testCase.selection)

			if testCase.err != nil {
				require.Error(t, err)
				assert.Equal(t, testCase.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.filtered, servers)
		})
	}
}
