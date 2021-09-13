package configuration

import (
	"errors"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/params/mock_params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Provider_readIvpn(t *testing.T) { //nolint:gocognit
	t.Parallel()

	var errDummy = errors.New("dummy test error")

	type singleStringCall struct {
		call  bool
		value string
		err   error
	}

	type portCall struct {
		getCall   bool
		getValue  string // "" or "0"
		getErr    error
		portCall  bool
		portValue uint16
		portErr   error
	}

	type sliceStringCall struct {
		call   bool
		values []string
		err    error
	}

	testCases := map[string]struct {
		targetIP  singleStringCall
		countries sliceStringCall
		cities    sliceStringCall
		isps      sliceStringCall
		hostnames sliceStringCall
		protocol  singleStringCall
		ovpnPort  portCall
		wgPort    portCall
		wgOldPort portCall
		settings  Provider
		err       error
	}{
		"target IP error": {
			targetIP: singleStringCall{call: true, value: "something", err: errDummy},
			settings: Provider{
				Name: constants.Ivpn,
			},
			err: errors.New("environment variable OPENVPN_TARGET_IP: dummy test error"),
		},
		"countries error": {
			targetIP:  singleStringCall{call: true},
			countries: sliceStringCall{call: true, err: errDummy},
			settings: Provider{
				Name: constants.Ivpn,
			},
			err: errors.New("environment variable COUNTRY: dummy test error"),
		},
		"cities error": {
			targetIP:  singleStringCall{call: true},
			countries: sliceStringCall{call: true},
			cities:    sliceStringCall{call: true, err: errDummy},
			settings: Provider{
				Name: constants.Ivpn,
			},
			err: errors.New("environment variable CITY: dummy test error"),
		},
		"isps error": {
			targetIP:  singleStringCall{call: true},
			countries: sliceStringCall{call: true},
			cities:    sliceStringCall{call: true},
			isps:      sliceStringCall{call: true, err: errDummy},
			settings: Provider{
				Name: constants.Ivpn,
			},
			err: errors.New("environment variable ISP: dummy test error"),
		},
		"hostnames error": {
			targetIP:  singleStringCall{call: true},
			countries: sliceStringCall{call: true},
			cities:    sliceStringCall{call: true},
			isps:      sliceStringCall{call: true},
			hostnames: sliceStringCall{call: true, err: errDummy},
			settings: Provider{
				Name: constants.Ivpn,
			},
			err: errors.New("environment variable SERVER_HOSTNAME: dummy test error"),
		},
		"openvpn protocol error": {
			targetIP:  singleStringCall{call: true},
			countries: sliceStringCall{call: true},
			cities:    sliceStringCall{call: true},
			isps:      sliceStringCall{call: true},
			hostnames: sliceStringCall{call: true},
			protocol:  singleStringCall{call: true, err: errDummy},
			settings: Provider{
				Name: constants.Ivpn,
			},
			err: errors.New("environment variable PROTOCOL: dummy test error"),
		},
		"openvpn custom port error": {
			targetIP:  singleStringCall{call: true},
			countries: sliceStringCall{call: true},
			cities:    sliceStringCall{call: true},
			isps:      sliceStringCall{call: true},
			hostnames: sliceStringCall{call: true},
			protocol:  singleStringCall{call: true},
			ovpnPort:  portCall{getCall: true, getErr: errDummy},
			settings: Provider{
				Name: constants.Ivpn,
			},
			err: errors.New("environment variable PORT: dummy test error"),
		},
		"wireguard custom port error": {
			targetIP:  singleStringCall{call: true},
			countries: sliceStringCall{call: true},
			cities:    sliceStringCall{call: true},
			isps:      sliceStringCall{call: true},
			hostnames: sliceStringCall{call: true},
			protocol:  singleStringCall{call: true},
			ovpnPort:  portCall{getCall: true, getValue: "0"},
			wgPort:    portCall{getCall: true, getErr: errDummy},
			settings: Provider{
				Name: constants.Ivpn,
			},
			err: errors.New("environment variable WIREGUARD_ENDPOINT_PORT: dummy test error"),
		},
		"default settings": {
			targetIP:  singleStringCall{call: true},
			countries: sliceStringCall{call: true},
			cities:    sliceStringCall{call: true},
			isps:      sliceStringCall{call: true},
			hostnames: sliceStringCall{call: true},
			protocol:  singleStringCall{call: true},
			ovpnPort:  portCall{getCall: true, getValue: "0"},
			wgPort:    portCall{getCall: true, getValue: "0"},
			wgOldPort: portCall{getCall: true, getValue: "0"},
			settings: Provider{
				Name: constants.Ivpn,
			},
		},
		"set settings": {
			targetIP:  singleStringCall{call: true, value: "1.2.3.4"},
			countries: sliceStringCall{call: true, values: []string{"A", "B"}},
			cities:    sliceStringCall{call: true, values: []string{"C", "D"}},
			isps:      sliceStringCall{call: true, values: []string{"ISP 1"}},
			hostnames: sliceStringCall{call: true, values: []string{"E", "F"}},
			protocol:  singleStringCall{call: true, value: constants.TCP},
			ovpnPort:  portCall{getCall: true, portCall: true, portValue: 443},
			wgPort:    portCall{getCall: true, portCall: true, portValue: 2049},
			settings: Provider{
				Name: constants.Ivpn,
				ServerSelection: ServerSelection{
					OpenVPN: OpenVPNSelection{
						TCP:        true,
						CustomPort: 443,
					},
					Wireguard: WireguardSelection{
						EndpointPort: 2049,
					},
					TargetIP:  net.IPv4(1, 2, 3, 4),
					Countries: []string{"A", "B"},
					Cities:    []string{"C", "D"},
					ISPs:      []string{"ISP 1"},
					Hostnames: []string{"E", "F"},
				},
			},
		},
	}
	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			env := mock_params.NewMockInterface(ctrl)

			servers := []models.IvpnServer{{Hostname: "a"}}
			allServers := models.AllServers{
				Ivpn: models.IvpnServers{
					Servers: servers,
				},
			}

			if testCase.targetIP.call {
				env.EXPECT().Get("OPENVPN_TARGET_IP").
					Return(testCase.targetIP.value, testCase.targetIP.err)
			}
			if testCase.countries.call {
				env.EXPECT().CSVInside("COUNTRY", constants.IvpnCountryChoices(servers)).
					Return(testCase.countries.values, testCase.countries.err)
			}
			if testCase.cities.call {
				env.EXPECT().CSVInside("CITY", constants.IvpnCityChoices(servers)).
					Return(testCase.cities.values, testCase.cities.err)
			}
			if testCase.isps.call {
				env.EXPECT().CSVInside("ISP", constants.IvpnISPChoices(servers)).
					Return(testCase.isps.values, testCase.isps.err)
			}
			if testCase.hostnames.call {
				env.EXPECT().CSVInside("SERVER_HOSTNAME", constants.IvpnHostnameChoices(servers)).
					Return(testCase.hostnames.values, testCase.hostnames.err)
			}
			if testCase.protocol.call {
				env.EXPECT().Inside("PROTOCOL", []string{constants.TCP, constants.UDP}, gomock.Any()).
					Return(testCase.protocol.value, testCase.protocol.err)
			}
			if testCase.ovpnPort.getCall {
				env.EXPECT().Get("PORT", gomock.Any()).
					Return(testCase.ovpnPort.getValue, testCase.ovpnPort.getErr)
			}
			if testCase.ovpnPort.portCall {
				env.EXPECT().Port("PORT").
					Return(testCase.ovpnPort.portValue, testCase.ovpnPort.portErr)
			}
			if testCase.wgPort.getCall {
				env.EXPECT().Get("WIREGUARD_ENDPOINT_PORT", gomock.Any()).
					Return(testCase.wgPort.getValue, testCase.wgPort.getErr)
			}
			if testCase.wgPort.portCall {
				env.EXPECT().Port("WIREGUARD_ENDPOINT_PORT").
					Return(testCase.wgPort.portValue, testCase.wgPort.portErr)
			}
			if testCase.wgOldPort.getCall {
				env.EXPECT().Get("WIREGUARD_PORT", gomock.Any()).
					Return(testCase.wgOldPort.getValue, testCase.wgOldPort.getErr)
			}
			if testCase.wgOldPort.portCall {
				env.EXPECT().Port("WIREGUARD_PORT").
					Return(testCase.wgOldPort.portValue, testCase.wgOldPort.portErr)
			}

			r := reader{
				servers: allServers,
				env:     env,
			}

			var settings Provider
			err := settings.readIvpn(r)

			if testCase.err != nil {
				require.Error(t, err)
				assert.Equal(t, testCase.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.settings, settings)
		})
	}
}
