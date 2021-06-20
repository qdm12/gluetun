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
					Group:   "group",
					Regions: []string{"a", "El country"},
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
		"fastestvpn": {
			settings: Provider{
				Name: constants.Fastestvpn,
				ServerSelection: ServerSelection{
					Hostnames: []string{"a", "b"},
					Countries: []string{"c", "d"},
				},
			},
			lines: []string{
				"|--Fastestvpn settings:",
				"   |--Network protocol: udp",
				"   |--Hostnames: a, b",
				"   |--Countries: c, d",
			},
		},
		"hidemyass": {
			settings: Provider{
				Name: constants.HideMyAss,
				ServerSelection: ServerSelection{
					Countries: []string{"a", "b"},
					Cities:    []string{"c", "d"},
					Hostnames: []string{"e", "f"},
				},
			},
			lines: []string{
				"|--Hidemyass settings:",
				"   |--Network protocol: udp",
				"   |--Countries: a, b",
				"   |--Cities: c, d",
				"   |--Hostnames: e, f",
			},
		},
		"ipvanish": {
			settings: Provider{
				Name: constants.Ipvanish,
				ServerSelection: ServerSelection{
					Countries: []string{"a", "b"},
					Cities:    []string{"c", "d"},
					Hostnames: []string{"e", "f"},
				},
			},
			lines: []string{
				"|--Ipvanish settings:",
				"   |--Network protocol: udp",
				"   |--Countries: a, b",
				"   |--Cities: c, d",
				"   |--Hostnames: e, f",
			},
		},
		"ivpn": {
			settings: Provider{
				Name: constants.Ivpn,
				ServerSelection: ServerSelection{
					Countries: []string{"a", "b"},
					Cities:    []string{"c", "d"},
					Hostnames: []string{"e", "f"},
				},
			},
			lines: []string{
				"|--Ivpn settings:",
				"   |--Network protocol: udp",
				"   |--Countries: a, b",
				"   |--Cities: c, d",
				"   |--Hostnames: e, f",
			},
		},
		"mullvad": {
			settings: Provider{
				Name: constants.Mullvad,
				ServerSelection: ServerSelection{
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
					Regions: []string{"a", "b"},
					Numbers: []uint16{1, 2},
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
					Hostnames: []string{"a", "b"},
				},
			},
			lines: []string{
				"|--Privado settings:",
				"   |--Network protocol: udp",
				"   |--Hostnames: a, b",
			},
		},
		"privatevpn": {
			settings: Provider{
				Name: constants.Privatevpn,
				ServerSelection: ServerSelection{
					Hostnames: []string{"a", "b"},
					Countries: []string{"c", "d"},
					Cities:    []string{"e", "f"},
				},
			},
			lines: []string{
				"|--Privatevpn settings:",
				"   |--Network protocol: udp",
				"   |--Countries: c, d",
				"   |--Cities: e, f",
				"   |--Hostnames: a, b",
			},
		},
		"protonvpn": {
			settings: Provider{
				Name: constants.Protonvpn,
				ServerSelection: ServerSelection{
					Countries: []string{"a", "b"},
					Regions:   []string{"c", "d"},
					Cities:    []string{"e", "f"},
					Names:     []string{"g", "h"},
					Hostnames: []string{"i", "j"},
				},
			},
			lines: []string{
				"|--Protonvpn settings:",
				"   |--Network protocol: udp",
				"   |--Countries: a, b",
				"   |--Regions: c, d",
				"   |--Cities: e, f",
				"   |--Names: g, h",
				"   |--Hostnames: i, j",
			},
		},
		"private internet access": {
			settings: Provider{
				Name: constants.PrivateInternetAccess,
				ServerSelection: ServerSelection{
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
					Regions: []string{"a", "b"},
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
		constants.VPNUnlimited: {
			settings: Provider{
				Name: constants.VPNUnlimited,
				ServerSelection: ServerSelection{
					Countries:  []string{"a", "b"},
					Cities:     []string{"c", "d"},
					Hostnames:  []string{"e", "f"},
					FreeOnly:   true,
					StreamOnly: true,
				},
				ExtraConfigOptions: ExtraConfigOptions{
					ClientKey: "a",
				},
			},
			lines: []string{
				"|--Vpn Unlimited settings:",
				"   |--Network protocol: udp",
				"   |--Countries: a, b",
				"   |--Cities: c, d",
				"   |--Hostnames: e, f",
				"   |--Free servers only",
				"   |--Stream servers only",
				"   |--Client key is set",
			},
		},
		"vyprvpn": {
			settings: Provider{
				Name: constants.Vyprvpn,
				ServerSelection: ServerSelection{
					Regions: []string{"a", "b"},
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
		mockStr string
		mockErr error
		tcp     bool
		err     error
	}{
		"error": {
			mockErr: errDummy,
			err:     errDummy,
		},
		"success": {
			mockStr: "tcp",
			tcp:     true,
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

			tcp, err := readProtocol(env)

			if testCase.err != nil {
				require.Error(t, err)
				assert.Equal(t, testCase.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.tcp, tcp)
		})
	}
}
