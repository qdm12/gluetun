package utils

import (
	"testing"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/constants/vpn"
	"github.com/stretchr/testify/assert"
)

func Test_getProtocol(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		selection settings.ServerSelection
		protocol  string
	}{
		"default": {
			protocol: constants.UDP,
		},
		"OpenVPN UDP": {
			selection: settings.ServerSelection{
				VPN: vpn.OpenVPN,
				OpenVPN: settings.OpenVPNSelection{
					Protocol: constants.UDP,
				},
			},
			protocol: constants.UDP,
		},
		"OpenVPN TCP": {
			selection: settings.ServerSelection{
				VPN: vpn.OpenVPN,
				OpenVPN: settings.OpenVPNSelection{
					Protocol: constants.TCP,
				},
			},
			protocol: constants.TCP,
		},
		"Wireguard": {
			selection: settings.ServerSelection{
				VPN: vpn.Wireguard,
			},
			protocol: constants.UDP,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			protocol := getProtocol(testCase.selection)

			assert.Equal(t, testCase.protocol, protocol)
		})
	}
}

func Test_filterByProtocol(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		selection settings.ServerSelection
		serverTCP bool
		serverUDP bool
		filtered  bool
	}{
		"Wireguard and server has UDP": {
			selection: settings.ServerSelection{
				VPN: vpn.Wireguard,
			},
			serverUDP: true,
			filtered:  false,
		},
		"Wireguard and server has not UDP": {
			selection: settings.ServerSelection{
				VPN: vpn.Wireguard,
			},
			serverUDP: false,
			filtered:  true,
		},
		"OpenVPN UDP and server has UDP": {
			selection: settings.ServerSelection{
				VPN: vpn.OpenVPN,
				OpenVPN: settings.OpenVPNSelection{
					Protocol: constants.UDP,
				},
			},
			serverUDP: true,
			filtered:  false,
		},
		"OpenVPN UDP and server has not UDP": {
			selection: settings.ServerSelection{
				VPN: vpn.OpenVPN,
				OpenVPN: settings.OpenVPNSelection{
					Protocol: constants.UDP,
				},
			},
			serverUDP: false,
			filtered:  true,
		},
		"OpenVPN TCP and server has TCP": {
			selection: settings.ServerSelection{
				VPN: vpn.OpenVPN,
				OpenVPN: settings.OpenVPNSelection{
					Protocol: constants.TCP,
				},
			},
			serverTCP: true,
			filtered:  false,
		},
		"OpenVPN TCP and server has not TCP": {
			selection: settings.ServerSelection{
				VPN: vpn.OpenVPN,
				OpenVPN: settings.OpenVPNSelection{
					Protocol: constants.TCP,
				},
			},
			serverTCP: false,
			filtered:  true,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			filtered := filterByProtocol(testCase.selection,
				testCase.serverTCP, testCase.serverUDP)

			assert.Equal(t, testCase.filtered, filtered)
		})
	}
}
