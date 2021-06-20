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

func Test_Provider_ipvanishLines(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		settings Provider
		lines    []string
	}{
		"empty settings": {},
		"full settings": {
			settings: Provider{
				ServerSelection: ServerSelection{
					Countries: []string{"A", "B"},
					Cities:    []string{"C", "D"},
					Hostnames: []string{"E", "F"},
				},
			},
			lines: []string{
				"|--Countries: A, B",
				"|--Cities: C, D",
				"|--Hostnames: E, F",
			},
		},
	}
	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			lines := testCase.settings.ipvanishLines()

			assert.Equal(t, testCase.lines, lines)
		})
	}
}

func Test_Provider_readIpvanish(t *testing.T) {
	t.Parallel()

	var errDummy = errors.New("dummy test error")

	type singleStringCall struct {
		call  bool
		value string
		err   error
	}

	type sliceStringCall struct {
		call   bool
		values []string
		err    error
	}

	testCases := map[string]struct {
		protocol  singleStringCall
		targetIP  singleStringCall
		countries sliceStringCall
		cities    sliceStringCall
		hostnames sliceStringCall
		settings  Provider
		err       error
	}{
		"protocol error": {
			protocol: singleStringCall{call: true, err: errDummy},
			settings: Provider{
				Name: constants.Ipvanish,
			},
			err: errDummy,
		},
		"target IP error": {
			protocol: singleStringCall{call: true},
			targetIP: singleStringCall{call: true, value: "something", err: errDummy},
			settings: Provider{
				Name: constants.Ipvanish,
			},
			err: errDummy,
		},
		"countries error": {
			protocol:  singleStringCall{call: true},
			targetIP:  singleStringCall{call: true},
			countries: sliceStringCall{call: true, err: errDummy},
			settings: Provider{
				Name: constants.Ipvanish,
			},
			err: errDummy,
		},
		"cities error": {
			protocol:  singleStringCall{call: true},
			targetIP:  singleStringCall{call: true},
			countries: sliceStringCall{call: true},
			cities:    sliceStringCall{call: true, err: errDummy},
			settings: Provider{
				Name: constants.Ipvanish,
			},
			err: errDummy,
		},
		"hostnames error": {
			protocol:  singleStringCall{call: true},
			targetIP:  singleStringCall{call: true},
			countries: sliceStringCall{call: true},
			cities:    sliceStringCall{call: true},
			hostnames: sliceStringCall{call: true, err: errDummy},
			settings: Provider{
				Name: constants.Ipvanish,
			},
			err: errDummy,
		},
		"default settings": {
			protocol:  singleStringCall{call: true},
			targetIP:  singleStringCall{call: true},
			countries: sliceStringCall{call: true},
			cities:    sliceStringCall{call: true},
			hostnames: sliceStringCall{call: true},
			settings: Provider{
				Name: constants.Ipvanish,
			},
		},
		"set settings": {
			protocol:  singleStringCall{call: true, value: constants.TCP},
			targetIP:  singleStringCall{call: true, value: "1.2.3.4"},
			countries: sliceStringCall{call: true, values: []string{"A", "B"}},
			cities:    sliceStringCall{call: true, values: []string{"C", "D"}},
			hostnames: sliceStringCall{call: true, values: []string{"E", "F"}},
			settings: Provider{
				Name: constants.Ipvanish,
				ServerSelection: ServerSelection{
					TCP:       true,
					TargetIP:  net.IPv4(1, 2, 3, 4),
					Countries: []string{"A", "B"},
					Cities:    []string{"C", "D"},
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

			env := mock_params.NewMockEnv(ctrl)
			if testCase.protocol.call {
				env.EXPECT().Inside("PROTOCOL", []string{constants.TCP, constants.UDP}, gomock.Any()).
					Return(testCase.protocol.value, testCase.protocol.err)
			}
			if testCase.targetIP.call {
				env.EXPECT().Get("OPENVPN_TARGET_IP").
					Return(testCase.targetIP.value, testCase.targetIP.err)
			}
			if testCase.countries.call {
				env.EXPECT().CSVInside("COUNTRY", constants.IpvanishCountryChoices()).
					Return(testCase.countries.values, testCase.countries.err)
			}
			if testCase.cities.call {
				env.EXPECT().CSVInside("CITY", constants.IpvanishCityChoices()).
					Return(testCase.cities.values, testCase.cities.err)
			}
			if testCase.hostnames.call {
				env.EXPECT().CSVInside("SERVER_HOSTNAME", constants.IpvanishHostnameChoices()).
					Return(testCase.hostnames.values, testCase.hostnames.err)
			}

			r := reader{env: env}

			var settings Provider
			err := settings.readIpvanish(r)

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
