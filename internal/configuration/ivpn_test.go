package configuration

import (
	"errors"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/golibs/params/mock_params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Provider_readIvpn(t *testing.T) {
	t.Parallel()

	var errDummy = errors.New("dummy test error")

	type singleStringCall struct {
		call  bool
		value string
		err   error
	}

	type singleUint16Call struct {
		call  bool
		value uint16
		err   error
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
		portGet   singleStringCall
		portPort  singleUint16Call
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
		"protocol error": {
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
		"custom port error": {
			targetIP:  singleStringCall{call: true},
			countries: sliceStringCall{call: true},
			cities:    sliceStringCall{call: true},
			isps:      sliceStringCall{call: true},
			hostnames: sliceStringCall{call: true},
			protocol:  singleStringCall{call: true},
			portGet:   singleStringCall{call: true, err: errDummy},
			settings: Provider{
				Name: constants.Ivpn,
			},
			err: errors.New("environment variable PORT: dummy test error"),
		},
		"default settings": {
			targetIP:  singleStringCall{call: true},
			countries: sliceStringCall{call: true},
			cities:    sliceStringCall{call: true},
			isps:      sliceStringCall{call: true},
			hostnames: sliceStringCall{call: true},
			protocol:  singleStringCall{call: true},
			portGet:   singleStringCall{call: true, value: "0"},
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
			portGet:   singleStringCall{call: true},
			portPort:  singleUint16Call{call: true, value: 443},
			settings: Provider{
				Name: constants.Ivpn,
				ServerSelection: ServerSelection{
					OpenVPN: OpenVPNSelection{
						TCP:        true,
						CustomPort: 443,
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
			if testCase.targetIP.call {
				env.EXPECT().Get("OPENVPN_TARGET_IP").
					Return(testCase.targetIP.value, testCase.targetIP.err)
			}
			if testCase.countries.call {
				env.EXPECT().CSVInside("COUNTRY", constants.IvpnCountryChoices()).
					Return(testCase.countries.values, testCase.countries.err)
			}
			if testCase.cities.call {
				env.EXPECT().CSVInside("CITY", constants.IvpnCityChoices()).
					Return(testCase.cities.values, testCase.cities.err)
			}
			if testCase.isps.call {
				env.EXPECT().CSVInside("ISP", constants.IvpnISPChoices()).
					Return(testCase.isps.values, testCase.isps.err)
			}
			if testCase.hostnames.call {
				env.EXPECT().CSVInside("SERVER_HOSTNAME", constants.IvpnHostnameChoices()).
					Return(testCase.hostnames.values, testCase.hostnames.err)
			}
			if testCase.protocol.call {
				env.EXPECT().Inside("PROTOCOL", []string{constants.TCP, constants.UDP}, gomock.Any()).
					Return(testCase.protocol.value, testCase.protocol.err)
			}
			if testCase.portGet.call {
				env.EXPECT().Get("PORT", gomock.Any()).
					Return(testCase.portGet.value, testCase.portGet.err)
			}
			if testCase.portPort.call {
				env.EXPECT().Port("PORT").
					Return(testCase.portPort.value, testCase.portPort.err)
			}

			r := reader{env: env}

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
