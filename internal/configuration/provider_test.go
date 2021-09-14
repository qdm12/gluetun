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
					VPN:       constants.OpenVPN,
					Countries: []string{"a", "El country"},
				},
			},
			lines: []string{
				"|--Cyberghost settings:",
				"   |--Countries: a, El country",
				"   |--OpenVPN selection:",
				"      |--Protocol: udp",
			},
		},
		"expressvpn": {
			settings: Provider{
				Name: constants.Expressvpn,
				ServerSelection: ServerSelection{
					VPN:       constants.OpenVPN,
					Hostnames: []string{"a", "b"},
					Countries: []string{"c", "d"},
					Cities:    []string{"e", "f"},
				},
			},
			lines: []string{
				"|--Expressvpn settings:",
				"   |--Countries: c, d",
				"   |--Cities: e, f",
				"   |--Hostnames: a, b",
				"   |--OpenVPN selection:",
				"      |--Protocol: udp",
			},
		},
		"fastestvpn": {
			settings: Provider{
				Name: constants.Fastestvpn,
				ServerSelection: ServerSelection{
					VPN:       constants.OpenVPN,
					Hostnames: []string{"a", "b"},
					Countries: []string{"c", "d"},
				},
			},
			lines: []string{
				"|--Fastestvpn settings:",
				"   |--Countries: c, d",
				"   |--Hostnames: a, b",
				"   |--OpenVPN selection:",
				"      |--Protocol: udp",
			},
		},
		"hidemyass": {
			settings: Provider{
				Name: constants.HideMyAss,
				ServerSelection: ServerSelection{
					VPN:       constants.OpenVPN,
					Countries: []string{"a", "b"},
					Cities:    []string{"c", "d"},
					Hostnames: []string{"e", "f"},
				},
			},
			lines: []string{
				"|--Hidemyass settings:",
				"   |--Countries: a, b",
				"   |--Cities: c, d",
				"   |--Hostnames: e, f",
				"   |--OpenVPN selection:",
				"      |--Protocol: udp",
			},
		},
		"ipvanish": {
			settings: Provider{
				Name: constants.Ipvanish,
				ServerSelection: ServerSelection{
					VPN:       constants.OpenVPN,
					Countries: []string{"a", "b"},
					Cities:    []string{"c", "d"},
					Hostnames: []string{"e", "f"},
				},
			},
			lines: []string{
				"|--Ipvanish settings:",
				"   |--Countries: a, b",
				"   |--Cities: c, d",
				"   |--Hostnames: e, f",
				"   |--OpenVPN selection:",
				"      |--Protocol: udp",
			},
		},
		"ivpn": {
			settings: Provider{
				Name: constants.Ivpn,
				ServerSelection: ServerSelection{
					VPN:       constants.OpenVPN,
					Countries: []string{"a", "b"},
					Cities:    []string{"c", "d"},
					Hostnames: []string{"e", "f"},
				},
			},
			lines: []string{
				"|--Ivpn settings:",
				"   |--Countries: a, b",
				"   |--Cities: c, d",
				"   |--Hostnames: e, f",
				"   |--OpenVPN selection:",
				"      |--Protocol: udp",
			},
		},
		"mullvad": {
			settings: Provider{
				Name: constants.Mullvad,
				ServerSelection: ServerSelection{
					VPN:       constants.OpenVPN,
					Countries: []string{"a", "b"},
					Cities:    []string{"c", "d"},
					ISPs:      []string{"e", "f"},
					OpenVPN: OpenVPNSelection{
						CustomPort: 1,
					},
				},
			},
			lines: []string{
				"|--Mullvad settings:",
				"   |--Countries: a, b",
				"   |--Cities: c, d",
				"   |--ISPs: e, f",
				"   |--OpenVPN selection:",
				"      |--Protocol: udp",
				"      |--Custom port: 1",
			},
		},
		"nordvpn": {
			settings: Provider{
				Name: constants.Nordvpn,
				ServerSelection: ServerSelection{
					VPN:     constants.OpenVPN,
					Regions: []string{"a", "b"},
					Numbers: []uint16{1, 2},
				},
			},
			lines: []string{
				"|--Nordvpn settings:",
				"   |--Regions: a, b",
				"   |--Numbers: 1, 2",
				"   |--OpenVPN selection:",
				"      |--Protocol: udp",
			},
		},
		"privado": {
			settings: Provider{
				Name: constants.Privado,
				ServerSelection: ServerSelection{
					VPN:       constants.OpenVPN,
					Hostnames: []string{"a", "b"},
				},
			},
			lines: []string{
				"|--Privado settings:",
				"   |--Hostnames: a, b",
				"   |--OpenVPN selection:",
				"      |--Protocol: udp",
			},
		},
		"privatevpn": {
			settings: Provider{
				Name: constants.Privatevpn,
				ServerSelection: ServerSelection{
					VPN:       constants.OpenVPN,
					Hostnames: []string{"a", "b"},
					Countries: []string{"c", "d"},
					Cities:    []string{"e", "f"},
				},
			},
			lines: []string{
				"|--Privatevpn settings:",
				"   |--Countries: c, d",
				"   |--Cities: e, f",
				"   |--Hostnames: a, b",
				"   |--OpenVPN selection:",
				"      |--Protocol: udp",
			},
		},
		"protonvpn": {
			settings: Provider{
				Name: constants.Protonvpn,
				ServerSelection: ServerSelection{
					VPN:       constants.OpenVPN,
					Countries: []string{"a", "b"},
					Regions:   []string{"c", "d"},
					Cities:    []string{"e", "f"},
					Names:     []string{"g", "h"},
					Hostnames: []string{"i", "j"},
				},
			},
			lines: []string{
				"|--Protonvpn settings:",
				"   |--Countries: a, b",
				"   |--Regions: c, d",
				"   |--Cities: e, f",
				"   |--Hostnames: i, j",
				"   |--Names: g, h",
				"   |--OpenVPN selection:",
				"      |--Protocol: udp",
			},
		},
		"private internet access": {
			settings: Provider{
				Name: constants.PrivateInternetAccess,
				ServerSelection: ServerSelection{
					VPN:     constants.OpenVPN,
					Regions: []string{"a", "b"},
					OpenVPN: OpenVPNSelection{
						CustomPort: 1,
					},
				},
				PortForwarding: PortForwarding{
					Enabled:  true,
					Filepath: string("/here"),
				},
			},
			lines: []string{
				"|--Private Internet Access settings:",
				"   |--Regions: a, b",
				"   |--OpenVPN selection:",
				"      |--Protocol: udp",
				"      |--Custom port: 1",
				"   |--Port forwarding:",
				"      |--File path: /here",
			},
		},
		"purevpn": {
			settings: Provider{
				Name: constants.Purevpn,
				ServerSelection: ServerSelection{
					VPN:       constants.OpenVPN,
					Regions:   []string{"a", "b"},
					Countries: []string{"c", "d"},
					Cities:    []string{"e", "f"},
				},
			},
			lines: []string{
				"|--Purevpn settings:",
				"   |--Countries: c, d",
				"   |--Regions: a, b",
				"   |--Cities: e, f",
				"   |--OpenVPN selection:",
				"      |--Protocol: udp",
			},
		},
		"surfshark": {
			settings: Provider{
				Name: constants.Surfshark,
				ServerSelection: ServerSelection{
					VPN:     constants.OpenVPN,
					Regions: []string{"a", "b"},
				},
			},
			lines: []string{
				"|--Surfshark settings:",
				"   |--Regions: a, b",
				"   |--OpenVPN selection:",
				"      |--Protocol: udp",
			},
		},
		"torguard": {
			settings: Provider{
				Name: constants.Torguard,
				ServerSelection: ServerSelection{
					VPN:       constants.OpenVPN,
					Countries: []string{"a", "b"},
					Cities:    []string{"c", "d"},
					Hostnames: []string{"e"},
				},
			},
			lines: []string{
				"|--Torguard settings:",
				"   |--Countries: a, b",
				"   |--Cities: c, d",
				"   |--Hostnames: e",
				"   |--OpenVPN selection:",
				"      |--Protocol: udp",
			},
		},
		constants.VPNUnlimited: {
			settings: Provider{
				Name: constants.VPNUnlimited,
				ServerSelection: ServerSelection{
					VPN:        constants.OpenVPN,
					Countries:  []string{"a", "b"},
					Cities:     []string{"c", "d"},
					Hostnames:  []string{"e", "f"},
					FreeOnly:   true,
					StreamOnly: true,
				},
			},
			lines: []string{
				"|--Vpn Unlimited settings:",
				"   |--Countries: a, b",
				"   |--Cities: c, d",
				"   |--Free servers only",
				"   |--Stream servers only",
				"   |--Hostnames: e, f",
				"   |--OpenVPN selection:",
				"      |--Protocol: udp",
			},
		},
		"vyprvpn": {
			settings: Provider{
				Name: constants.Vyprvpn,
				ServerSelection: ServerSelection{
					VPN:     constants.OpenVPN,
					Regions: []string{"a", "b"},
				},
			},
			lines: []string{
				"|--Vyprvpn settings:",
				"   |--Regions: a, b",
				"   |--OpenVPN selection:",
				"      |--Protocol: udp",
			},
		},
		"wevpn": {
			settings: Provider{
				Name: constants.Wevpn,
				ServerSelection: ServerSelection{
					VPN:       constants.OpenVPN,
					Cities:    []string{"a", "b"},
					Hostnames: []string{"c", "d"},
					OpenVPN: OpenVPNSelection{
						CustomPort: 1,
					},
				},
			},
			lines: []string{
				"|--Wevpn settings:",
				"   |--Cities: a, b",
				"   |--Hostnames: c, d",
				"   |--OpenVPN selection:",
				"      |--Protocol: udp",
				"      |--Custom port: 1",
			},
		},
		"windscribe": {
			settings: Provider{
				Name: constants.Windscribe,
				ServerSelection: ServerSelection{
					VPN:       constants.OpenVPN,
					Regions:   []string{"a", "b"},
					Cities:    []string{"c", "d"},
					Hostnames: []string{"e", "f"},
					OpenVPN: OpenVPNSelection{
						CustomPort: 1,
					},
				},
			},
			lines: []string{
				"|--Windscribe settings:",
				"   |--Regions: a, b",
				"   |--Cities: c, d",
				"   |--Hostnames: e, f",
				"   |--OpenVPN selection:",
				"      |--Protocol: udp",
				"      |--Custom port: 1",
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
			err:     errors.New("environment variable OPENVPN_PROTOCOL: dummy"),
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

			env := mock_params.NewMockInterface(ctrl)
			env.EXPECT().
				Inside("OPENVPN_PROTOCOL", []string{"tcp", "udp"}, gomock.Any(), gomock.Any()).
				Return(testCase.mockStr, testCase.mockErr)
			reader := reader{
				env: env,
			}

			tcp, err := readOpenVPNProtocol(reader)

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
