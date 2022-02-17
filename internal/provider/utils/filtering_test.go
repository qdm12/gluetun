package utils

import (
	"testing"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/constants/vpn"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/stretchr/testify/assert"
)

func Test_FilterServers(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		servers   []models.Server
		selection settings.ServerSelection
		filtered  []models.Server
	}{
		"no server available": {
			selection: settings.ServerSelection{}.WithDefaults(providers.Mullvad),
		},
		"no filter": {
			servers: []models.Server{
				{VPN: vpn.OpenVPN, Hostname: "a", UDP: true},
				{VPN: vpn.OpenVPN, Hostname: "b", UDP: true},
				{VPN: vpn.OpenVPN, Hostname: "c", UDP: true},
			},
			selection: settings.ServerSelection{}.WithDefaults(providers.Mullvad),
			filtered: []models.Server{
				{VPN: vpn.OpenVPN, Hostname: "a", UDP: true},
				{VPN: vpn.OpenVPN, Hostname: "b", UDP: true},
				{VPN: vpn.OpenVPN, Hostname: "c", UDP: true},
			},
		},
		"filter by VPN protocol": {
			selection: settings.ServerSelection{
				VPN: vpn.Wireguard,
			}.WithDefaults(providers.Mullvad),
			servers: []models.Server{
				{VPN: vpn.OpenVPN, Hostname: "a", UDP: true},
				{VPN: vpn.Wireguard, Hostname: "b", UDP: true},
				{VPN: vpn.OpenVPN, Hostname: "c", UDP: true},
			},
			filtered: []models.Server{
				{VPN: vpn.Wireguard, Hostname: "b", UDP: true},
			},
		},
		"filter by network protocol": {
			selection: settings.ServerSelection{
				OpenVPN: settings.OpenVPNSelection{
					TCP: boolPtr(true),
				},
			}.WithDefaults(providers.Ivpn),
			servers: []models.Server{
				{UDP: true, Hostname: "a", VPN: vpn.OpenVPN},
				{UDP: true, TCP: true, Hostname: "b", VPN: vpn.OpenVPN},
				{UDP: true, Hostname: "c", VPN: vpn.OpenVPN},
			},
			filtered: []models.Server{
				{UDP: true, TCP: true, Hostname: "b", VPN: vpn.OpenVPN},
			},
		},
		"filter by multihop only": {
			selection: settings.ServerSelection{
				MultiHopOnly: boolPtr(true),
			}.WithDefaults(providers.Surfshark),
			servers: []models.Server{
				{MultiHop: false, VPN: vpn.OpenVPN, UDP: true},
				{MultiHop: true, VPN: vpn.OpenVPN, UDP: true},
				{MultiHop: false, VPN: vpn.OpenVPN, UDP: true},
			},
			filtered: []models.Server{
				{MultiHop: true, VPN: vpn.OpenVPN, UDP: true},
			},
		},
		"filter by free only": {
			selection: settings.ServerSelection{
				FreeOnly: boolPtr(true),
			}.WithDefaults(providers.Surfshark),
			servers: []models.Server{
				{Free: false, VPN: vpn.OpenVPN, UDP: true},
				{Free: true, VPN: vpn.OpenVPN, UDP: true},
				{Free: false, VPN: vpn.OpenVPN, UDP: true},
			},
			filtered: []models.Server{
				{Free: true, VPN: vpn.OpenVPN, UDP: true},
			},
		},
		"filter by premium only": {
			selection: settings.ServerSelection{
				PremiumOnly: boolPtr(true),
			}.WithDefaults(providers.Surfshark),
			servers: []models.Server{
				{Premium: false, VPN: vpn.OpenVPN, UDP: true},
				{Premium: true, VPN: vpn.OpenVPN, UDP: true},
				{Premium: false, VPN: vpn.OpenVPN, UDP: true},
			},
			filtered: []models.Server{
				{Premium: true, VPN: vpn.OpenVPN, UDP: true},
			},
		},
		"filter by stream only": {
			selection: settings.ServerSelection{
				StreamOnly: boolPtr(true),
			}.WithDefaults(providers.Surfshark),
			servers: []models.Server{
				{Stream: false, VPN: vpn.OpenVPN, UDP: true},
				{Stream: true, VPN: vpn.OpenVPN, UDP: true},
				{Stream: false, VPN: vpn.OpenVPN, UDP: true},
			},
			filtered: []models.Server{
				{Stream: true, VPN: vpn.OpenVPN, UDP: true},
			},
		},
		"filter by owned": {
			selection: settings.ServerSelection{
				OwnedOnly: boolPtr(true),
			}.WithDefaults(providers.Mullvad),
			servers: []models.Server{
				{Owned: false, VPN: vpn.OpenVPN, UDP: true},
				{Owned: true, VPN: vpn.OpenVPN, UDP: true},
				{Owned: false, VPN: vpn.OpenVPN, UDP: true},
			},
			filtered: []models.Server{
				{Owned: true, VPN: vpn.OpenVPN, UDP: true},
			},
		},
		"filter by country": {
			selection: settings.ServerSelection{
				Countries: []string{"b"},
			}.WithDefaults(providers.Mullvad),
			servers: []models.Server{
				{Country: "a", VPN: vpn.OpenVPN, UDP: true},
				{Country: "b", VPN: vpn.OpenVPN, UDP: true},
				{Country: "c", VPN: vpn.OpenVPN, UDP: true},
			},
			filtered: []models.Server{
				{Country: "b", VPN: vpn.OpenVPN, UDP: true},
			},
		},
		"filter by region": {
			selection: settings.ServerSelection{
				Regions: []string{"b"},
			}.WithDefaults(providers.Surfshark),
			servers: []models.Server{
				{Region: "a", VPN: vpn.OpenVPN, UDP: true},
				{Region: "b", VPN: vpn.OpenVPN, UDP: true},
				{Region: "c", VPN: vpn.OpenVPN, UDP: true},
			},
			filtered: []models.Server{
				{Region: "b", VPN: vpn.OpenVPN, UDP: true},
			},
		},
		"filter by city": {
			selection: settings.ServerSelection{
				Cities: []string{"b"},
			}.WithDefaults(providers.Mullvad),
			servers: []models.Server{
				{City: "a", VPN: vpn.OpenVPN, UDP: true},
				{City: "b", VPN: vpn.OpenVPN, UDP: true},
				{City: "c", VPN: vpn.OpenVPN, UDP: true},
			},
			filtered: []models.Server{
				{City: "b", VPN: vpn.OpenVPN, UDP: true},
			},
		},
		"filter by ISP": {
			selection: settings.ServerSelection{
				ISPs: []string{"b"},
			}.WithDefaults(providers.Mullvad),
			servers: []models.Server{
				{ISP: "a", VPN: vpn.OpenVPN, UDP: true},
				{ISP: "b", VPN: vpn.OpenVPN, UDP: true},
				{ISP: "c", VPN: vpn.OpenVPN, UDP: true},
			},
			filtered: []models.Server{
				{ISP: "b", VPN: vpn.OpenVPN, UDP: true},
			},
		},
		"filter by number": {
			selection: settings.ServerSelection{
				Numbers: []uint16{1},
			}.WithDefaults(providers.Mullvad),
			servers: []models.Server{
				{Number: 0, VPN: vpn.OpenVPN, UDP: true},
				{Number: 1, VPN: vpn.OpenVPN, UDP: true},
				{Number: 2, VPN: vpn.OpenVPN, UDP: true},
			},
			filtered: []models.Server{
				{Number: 1, VPN: vpn.OpenVPN, UDP: true},
			},
		},
		"filter by server name": {
			selection: settings.ServerSelection{
				Names: []string{"b"},
			}.WithDefaults(providers.Mullvad),
			servers: []models.Server{
				{ServerName: "a", VPN: vpn.OpenVPN, UDP: true},
				{ServerName: "b", VPN: vpn.OpenVPN, UDP: true},
				{ServerName: "c", VPN: vpn.OpenVPN, UDP: true},
			},
			filtered: []models.Server{
				{ServerName: "b", VPN: vpn.OpenVPN, UDP: true},
			},
		},
		"filter by hostname": {
			selection: settings.ServerSelection{
				Hostnames: []string{"b"},
			}.WithDefaults(providers.Mullvad),
			servers: []models.Server{
				{Hostname: "a", VPN: vpn.OpenVPN, UDP: true},
				{Hostname: "b", VPN: vpn.OpenVPN, UDP: true},
				{Hostname: "c", VPN: vpn.OpenVPN, UDP: true},
			},
			filtered: []models.Server{
				{Hostname: "b", VPN: vpn.OpenVPN, UDP: true},
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			filtered := filterServers(testCase.servers, testCase.selection)

			assert.Equal(t, testCase.filtered, filtered)
		})
	}
}

func Test_filterByPossibilities(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		value         string
		possibilities []string
		filtered      bool
	}{
		"no possibilities": {},
		"value not in possibilities": {
			value:         "c",
			possibilities: []string{"a", "b"},
			filtered:      true,
		},
		"value in possibilities": {
			value:         "c",
			possibilities: []string{"a", "b", "c"},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			filtered := filterByPossibilities(testCase.value, testCase.possibilities)
			assert.Equal(t, testCase.filtered, filtered)
		})
	}
}
