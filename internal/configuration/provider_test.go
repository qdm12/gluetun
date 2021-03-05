package configuration

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/golibs/params/mock_params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errDummy = errors.New("dummy")

func Test_Provider_lines(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		settings Provider
		lines    []string
	}{
		"cyberghost": {
			settings: Provider{
				Name: constants.Cyberghost,
				ServerSelection: ServerSelection{
					Protocol: constants.UDP,
					Group:    "group",
					Regions:  []string{"a", "El country"},
				},
				ExtraConfigOptions: ExtraConfigOptions{
					ClientKey:         "a",
					ClientCertificate: "a",
				},
			},
			lines: []string{
				"|--Cyberghost settings:",
				"   |--Network protocol: udp",
				"   |--Server group: group",
				"   |--Regions: a, El country",
				"   |--Client key is set",
				"   |--Client certificate is set",
			},
		},
		"hidemyass": {
			settings: Provider{
				Name: constants.HideMyAss,
				ServerSelection: ServerSelection{
					Protocol:  constants.UDP,
					Countries: []string{"a", "b"},
					Cities:    []string{"c", "d"},
				},
			},
			lines: []string{
				"|--HideMyAss settings:",
				"   |--Network protocol: udp",
				"   |--Countries: a, b",
				"   |--Cities: c, d",
			},
		},
		"mullvad": {
			settings: Provider{
				Name: constants.Mullvad,
				ServerSelection: ServerSelection{
					Protocol:   constants.UDP,
					Countries:  []string{"a", "b"},
					Cities:     []string{"c", "d"},
					ISPs:       []string{"e", "f"},
					CustomPort: 1,
				},
				ExtraConfigOptions: ExtraConfigOptions{
					OpenVPNIPv6: true,
				},
			},
			lines: []string{
				"|--Mullvad settings:",
				"   |--Network protocol: udp",
				"   |--Countries: a, b",
				"   |--Cities: c, d",
				"   |--ISPs: e, f",
				"   |--Custom port: 1",
				"   |--IPv6: enabled",
			},
		},
		"nordvpn": {
			settings: Provider{
				Name: constants.Nordvpn,
				ServerSelection: ServerSelection{
					Protocol: constants.UDP,
					Regions:  []string{"a", "b"},
					Numbers:  []uint16{1, 2},
				},
			},
			lines: []string{
				"|--Nordvpn settings:",
				"   |--Network protocol: udp",
				"   |--Regions: a, b",
				"   |--Numbers: 1, 2",
			},
		},
		"privado": {
			settings: Provider{
				Name: constants.Privado,
				ServerSelection: ServerSelection{
					Protocol:  constants.UDP,
					Hostnames: []string{"a", "b"},
				},
			},
			lines: []string{
				"|--Privado settings:",
				"   |--Network protocol: udp",
				"   |--Hostnames: a, b",
			},
		},
		"private internet access": {
			settings: Provider{
				Name: constants.PrivateInternetAccess,
				ServerSelection: ServerSelection{
					Protocol:         constants.UDP,
					Regions:          []string{"a", "b"},
					EncryptionPreset: constants.PIAEncryptionPresetStrong,
					CustomPort:       1,
				},
				PortForwarding: PortForwarding{
					Enabled:  true,
					Filepath: string("/here"),
				},
			},
			lines: []string{
				"|--Private Internet Access settings:",
				"   |--Network protocol: udp",
				"   |--Regions: a, b",
				"   |--Encryption preset: strong",
				"   |--Custom port: 1",
				"   |--Port forwarding:",
				"      |--File path: /here",
			},
		},
		"purevpn": {
			settings: Provider{
				Name: constants.Purevpn,
				ServerSelection: ServerSelection{
					Protocol:  constants.UDP,
					Regions:   []string{"a", "b"},
					Countries: []string{"c", "d"},
					Cities:    []string{"e", "f"},
				},
			},
			lines: []string{
				"|--Purevpn settings:",
				"   |--Network protocol: udp",
				"   |--Regions: a, b",
				"   |--Countries: c, d",
				"   |--Cities: e, f",
			},
		},
		"surfshark": {
			settings: Provider{
				Name: constants.Surfshark,
				ServerSelection: ServerSelection{
					Protocol: constants.UDP,
					Regions:  []string{"a", "b"},
				},
			},
			lines: []string{
				"|--Surfshark settings:",
				"   |--Network protocol: udp",
				"   |--Regions: a, b",
			},
		},
		"torguard": {
			settings: Provider{
				Name: constants.Torguard,
				ServerSelection: ServerSelection{
					Protocol:  constants.UDP,
					Countries: []string{"a", "b"},
					Cities:    []string{"c", "d"},
					Hostnames: []string{"e"},
				},
			},
			lines: []string{
				"|--Torguard settings:",
				"   |--Network protocol: udp",
				"   |--Countries: a, b",
				"   |--Cities: c, d",
				"   |--Hostnames: e",
			},
		},
		"vyprvpn": {
			settings: Provider{
				Name: constants.Vyprvpn,
				ServerSelection: ServerSelection{
					Protocol: constants.UDP,
					Regions:  []string{"a", "b"},
				},
			},
			lines: []string{
				"|--Vyprvpn settings:",
				"   |--Network protocol: udp",
				"   |--Regions: a, b",
			},
		},
		"windscribe": {
			settings: Provider{
				Name: constants.Windscribe,
				ServerSelection: ServerSelection{
					Protocol:   constants.UDP,
					Regions:    []string{"a", "b"},
					Cities:     []string{"c", "d"},
					Hostnames:  []string{"e", "f"},
					CustomPort: 1,
				},
			},
			lines: []string{
				"|--Windscribe settings:",
				"   |--Network protocol: udp",
				"   |--Regions: a, b",
				"   |--Cities: c, d",
				"   |--Hostnames: e, f",
				"   |--Custom port: 1",
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			lines := testCase.settings.lines()

			assert.Equal(t, testCase.lines, lines)
		})
	}
}

func Test_readProtocol(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		mockStr  string
		mockErr  error
		protocol string
		err      error
	}{
		"error": {
			mockErr: errDummy,
			err:     errDummy,
		},
		"success": {
			mockStr:  "tcp",
			protocol: constants.TCP,
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			env := mock_params.NewMockEnv(ctrl)
			env.EXPECT().
				Inside("PROTOCOL", []string{"tcp", "udp"}, gomock.Any()).
				Return(testCase.mockStr, testCase.mockErr)

			protocol, err := readProtocol(env)

			if testCase.err != nil {
				require.Error(t, err)
				assert.Equal(t, testCase.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.protocol, protocol)
		})
	}
}
