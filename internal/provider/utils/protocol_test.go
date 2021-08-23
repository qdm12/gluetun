package utils

import (
	"testing"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/stretchr/testify/assert"
)

func Test_GetProtocol(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		selection configuration.ServerSelection
		protocol  string
	}{
		"default": {
			protocol: constants.UDP,
		},
		"OpenVPN UDP": {
			selection: configuration.ServerSelection{
				VPN: constants.OpenVPN,
			},
			protocol: constants.UDP,
		},
		"OpenVPN TCP": {
			selection: configuration.ServerSelection{
				VPN: constants.OpenVPN,
				OpenVPN: configuration.OpenVPNSelection{
					TCP: true,
				},
			},
			protocol: constants.TCP,
		},
		"Wireguard": {
			selection: configuration.ServerSelection{
				VPN: constants.Wireguard,
			},
			protocol: constants.UDP,
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			protocol := GetProtocol(testCase.selection)

			assert.Equal(t, testCase.protocol, protocol)
		})
	}
}

func Test_FilterByProtocol(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		selection configuration.ServerSelection
		serverTCP bool
		serverUDP bool
		filtered  bool
	}{
		"Wireguard and server has UDP": {
			selection: configuration.ServerSelection{
				VPN: constants.Wireguard,
			},
			serverUDP: true,
			filtered:  false,
		},
		"Wireguard and server has not UDP": {
			selection: configuration.ServerSelection{
				VPN: constants.Wireguard,
			},
			serverUDP: false,
			filtered:  true,
		},
		"OpenVPN UDP and server has UDP": {
			selection: configuration.ServerSelection{
				VPN: constants.OpenVPN,
				OpenVPN: configuration.OpenVPNSelection{
					TCP: false,
				},
			},
			serverUDP: true,
			filtered:  false,
		},
		"OpenVPN UDP and server has not UDP": {
			selection: configuration.ServerSelection{
				VPN: constants.OpenVPN,
				OpenVPN: configuration.OpenVPNSelection{
					TCP: false,
				},
			},
			serverUDP: false,
			filtered:  true,
		},
		"OpenVPN TCP and server has TCP": {
			selection: configuration.ServerSelection{
				VPN: constants.OpenVPN,
				OpenVPN: configuration.OpenVPNSelection{
					TCP: true,
				},
			},
			serverTCP: true,
			filtered:  false,
		},
		"OpenVPN TCP and server has not TCP": {
			selection: configuration.ServerSelection{
				VPN: constants.OpenVPN,
				OpenVPN: configuration.OpenVPNSelection{
					TCP: true,
				},
			},
			serverTCP: false,
			filtered:  true,
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			filtered := FilterByProtocol(testCase.selection,
				testCase.serverTCP, testCase.serverUDP)

			assert.Equal(t, testCase.filtered, filtered)
		})
	}
}
