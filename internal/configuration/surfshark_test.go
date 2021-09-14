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

func Test_Provider_readSurfshark(t *testing.T) {
	t.Parallel()

	var errDummy = errors.New("dummy test error")

	type stringCall struct {
		call  bool
		value string
		err   error
	}

	type boolCall struct {
		call  bool
		value bool
		err   error
	}

	type sliceStringCall struct {
		call   bool
		values []string
		err    error
	}

	testCases := map[string]struct {
		targetIP  stringCall
		countries sliceStringCall
		cities    sliceStringCall
		hostnames sliceStringCall
		regions   sliceStringCall
		multiHop  boolCall
		protocol  stringCall
		settings  Provider
		err       error
	}{
		"target IP error": {
			targetIP: stringCall{call: true, value: "something", err: errDummy},
			settings: Provider{
				Name: constants.Surfshark,
			},
			err: errors.New("environment variable OPENVPN_TARGET_IP: dummy test error"),
		},
		"countries error": {
			targetIP:  stringCall{call: true},
			countries: sliceStringCall{call: true, err: errDummy},
			settings: Provider{
				Name: constants.Surfshark,
			},
			err: errors.New("environment variable COUNTRY: dummy test error"),
		},
		"cities error": {
			targetIP:  stringCall{call: true},
			countries: sliceStringCall{call: true},
			cities:    sliceStringCall{call: true, err: errDummy},
			settings: Provider{
				Name: constants.Surfshark,
			},
			err: errors.New("environment variable CITY: dummy test error"),
		},
		"hostnames error": {
			targetIP:  stringCall{call: true},
			countries: sliceStringCall{call: true},
			cities:    sliceStringCall{call: true},
			hostnames: sliceStringCall{call: true, err: errDummy},
			settings: Provider{
				Name: constants.Surfshark,
			},
			err: errors.New("environment variable SERVER_HOSTNAME: dummy test error"),
		},
		"regions error": {
			targetIP:  stringCall{call: true},
			countries: sliceStringCall{call: true},
			cities:    sliceStringCall{call: true},
			hostnames: sliceStringCall{call: true},
			regions:   sliceStringCall{call: true, err: errDummy},
			settings: Provider{
				Name: constants.Surfshark,
			},
			err: errors.New("environment variable REGION: dummy test error"),
		},
		"multi hop error": {
			targetIP:  stringCall{call: true},
			countries: sliceStringCall{call: true},
			cities:    sliceStringCall{call: true},
			hostnames: sliceStringCall{call: true},
			regions:   sliceStringCall{call: true},
			multiHop:  boolCall{call: true, err: errDummy},
			settings: Provider{
				Name: constants.Surfshark,
			},
			err: errors.New("environment variable MULTIHOP_ONLY: dummy test error"),
		},
		"openvpn protocol error": {
			targetIP:  stringCall{call: true},
			countries: sliceStringCall{call: true},
			cities:    sliceStringCall{call: true},
			hostnames: sliceStringCall{call: true},
			regions:   sliceStringCall{call: true},
			multiHop:  boolCall{call: true},
			protocol:  stringCall{call: true, err: errDummy},
			settings: Provider{
				Name: constants.Surfshark,
			},
			err: errors.New("environment variable OPENVPN_PROTOCOL: dummy test error"),
		},
		"default settings": {
			targetIP:  stringCall{call: true},
			countries: sliceStringCall{call: true},
			cities:    sliceStringCall{call: true},
			hostnames: sliceStringCall{call: true},
			regions:   sliceStringCall{call: true},
			multiHop:  boolCall{call: true},
			protocol:  stringCall{call: true},
			settings: Provider{
				Name: constants.Surfshark,
			},
		},
		"set settings": {
			targetIP:  stringCall{call: true, value: "1.2.3.4"},
			countries: sliceStringCall{call: true, values: []string{"A", "B"}},
			cities:    sliceStringCall{call: true, values: []string{"C", "D"}},
			regions: sliceStringCall{call: true, values: []string{
				"E", "F", "netherlands amsterdam",
			}}, // Netherlands Amsterdam is for retro compatibility test
			multiHop:  boolCall{call: true},
			hostnames: sliceStringCall{call: true, values: []string{"E", "F"}},
			protocol:  stringCall{call: true, value: constants.TCP},
			settings: Provider{
				Name: constants.Surfshark,
				ServerSelection: ServerSelection{
					OpenVPN: OpenVPNSelection{
						TCP: true,
					},
					TargetIP:  net.IPv4(1, 2, 3, 4),
					Regions:   []string{"E", "F", "europe"},
					Countries: []string{"A", "B", "netherlands"},
					Cities:    []string{"C", "D", "amsterdam"},
					Hostnames: []string{"E", "F", "nl-ams.prod.surfshark.com"},
				},
			},
		},
		"Netherlands Amsterdam": {
			targetIP:  stringCall{call: true},
			countries: sliceStringCall{call: true},
			cities:    sliceStringCall{call: true},
			regions:   sliceStringCall{call: true, values: []string{"netherlands amsterdam"}},
			multiHop:  boolCall{call: true},
			hostnames: sliceStringCall{call: true},
			protocol:  stringCall{call: true},
			settings: Provider{
				Name: constants.Surfshark,
				ServerSelection: ServerSelection{
					Regions:   []string{"europe"},
					Countries: []string{"netherlands"},
					Cities:    []string{"amsterdam"},
					Hostnames: []string{"nl-ams.prod.surfshark.com"},
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

			servers := []models.SurfsharkServer{{Hostname: "a"}}
			allServers := models.AllServers{
				Surfshark: models.SurfsharkServers{
					Servers: servers,
				},
			}

			if testCase.targetIP.call {
				env.EXPECT().Get("OPENVPN_TARGET_IP").
					Return(testCase.targetIP.value, testCase.targetIP.err)
			}
			if testCase.countries.call {
				env.EXPECT().CSVInside("COUNTRY", constants.SurfsharkCountryChoices(servers)).
					Return(testCase.countries.values, testCase.countries.err)
			}
			if testCase.cities.call {
				env.EXPECT().CSVInside("CITY", constants.SurfsharkCityChoices(servers)).
					Return(testCase.cities.values, testCase.cities.err)
			}
			if testCase.hostnames.call {
				env.EXPECT().CSVInside("SERVER_HOSTNAME", constants.SurfsharkHostnameChoices(servers)).
					Return(testCase.hostnames.values, testCase.hostnames.err)
			}
			if testCase.regions.call {
				regionChoices := constants.SurfsharkRegionChoices(servers)
				regionChoices = append(regionChoices, constants.SurfsharkRetroLocChoices(servers)...)
				env.EXPECT().CSVInside("REGION", regionChoices).
					Return(testCase.regions.values, testCase.regions.err)
			}
			if testCase.multiHop.call {
				env.EXPECT().YesNo("MULTIHOP_ONLY", gomock.Any()).
					Return(testCase.multiHop.value, testCase.multiHop.err)
			}
			if testCase.protocol.call {
				env.EXPECT().Inside("OPENVPN_PROTOCOL", []string{constants.TCP, constants.UDP}, gomock.Any()).
					Return(testCase.protocol.value, testCase.protocol.err)
			}

			r := reader{
				servers: allServers,
				env:     env,
			}

			var settings Provider
			err := settings.readSurfshark(r)

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

func Test_surfsharkRetroRegion(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		original ServerSelection
		modified ServerSelection
	}{
		"empty": {},
		"1 retro region": {
			original: ServerSelection{
				Regions: []string{"australia adelaide"},
			},
			modified: ServerSelection{
				Regions:   []string{"asia pacific"},
				Countries: []string{"australia"},
				Cities:    []string{"adelaide"},
				Hostnames: []string{"au-adl.prod.surfshark.com"},
			},
		},
		"2 overlapping retro regions": {
			original: ServerSelection{
				Regions: []string{"australia adelaide", "australia melbourne"},
			},
			modified: ServerSelection{
				Regions:   []string{"asia pacific"},
				Countries: []string{"australia"},
				Cities:    []string{"adelaide", "melbourne"},
				Hostnames: []string{"au-adl.prod.surfshark.com", "au-mel.prod.surfshark.com"},
			},
		},
		"2 distinct retro regions": {
			original: ServerSelection{
				Regions: []string{"australia adelaide", "netherlands amsterdam"},
			},
			modified: ServerSelection{
				Regions:   []string{"asia pacific", "europe"},
				Countries: []string{"australia", "netherlands"},
				Cities:    []string{"adelaide", "amsterdam"},
				Hostnames: []string{"au-adl.prod.surfshark.com", "nl-ams.prod.surfshark.com"},
			},
		},
		"retro region with existing region": {
			// note TestRegion will be ignored in the filters downstream
			original: ServerSelection{
				Regions: []string{"TestRegion", "australia adelaide"},
			},
			modified: ServerSelection{
				Regions:   []string{"TestRegion", "asia pacific"},
				Countries: []string{"australia"},
				Cities:    []string{"adelaide"},
				Hostnames: []string{"au-adl.prod.surfshark.com"},
			},
		},
	}
	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			selection := surfsharkRetroRegion(testCase.original)

			assert.Equal(t, testCase.modified, selection)
		})
	}
}
