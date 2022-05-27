package storage

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func populateProviders(allProviderVersion uint16, allProviderTimestamp int64,
	servers models.AllServers) models.AllServers {
	allProviders := providers.All()
	if servers.ProviderToServers == nil {
		servers.ProviderToServers = make(map[string]models.Servers, len(allProviders)-1)
	}
	for _, provider := range allProviders {
		_, has := servers.ProviderToServers[provider]
		if has {
			continue
		}
		servers.ProviderToServers[provider] = models.Servers{
			Version:   allProviderVersion,
			Timestamp: allProviderTimestamp,
		}
	}
	return servers
}

func Test_extractServersFromBytes(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		b          []byte
		hardcoded  models.AllServers
		logged     []string
		persisted  models.AllServers
		errMessage string
	}{
		"bad JSON": {
			b:          []byte("garbage"),
			errMessage: "cannot decode servers: invalid character 'g' looking for beginning of value",
		},
		"bad provider JSON": {
			b:         []byte(`{"cyberghost": "garbage"}`),
			hardcoded: populateProviders(1, 0, models.AllServers{}),
			errMessage: "cannot decode servers for provider: Cyberghost: " +
				"json: cannot unmarshal string into Go value of type models.Servers",
		},
		"absent provider keys": {
			b:         []byte(`{}`),
			hardcoded: populateProviders(1, 0, models.AllServers{}),
			persisted: models.AllServers{
				ProviderToServers: map[string]models.Servers{},
			},
		},
		"same versions": {
			b: []byte(`{
					"cyberghost": {"version": 1, "timestamp": 1},
					"expressvpn": {"version": 1, "timestamp": 1},
					"fastestvpn": {"version": 1, "timestamp": 1},
					"hidemyass": {"version": 1, "timestamp": 1},
					"ipvanish": {"version": 1, "timestamp": 1},
					"ivpn": {"version": 1, "timestamp": 1},
					"mullvad": {"version": 1, "timestamp": 1},
					"nordvpn": {"version": 1, "timestamp": 1},
					"perfect privacy": {"version": 1, "timestamp": 1},
					"privado": {"version": 1, "timestamp": 1},
					"private internet access": {"version": 1, "timestamp": 1},
					"privatevpn": {"version": 1, "timestamp": 1},
					"protonvpn": {"version": 1, "timestamp": 1},
					"purevpn": {"version": 1, "timestamp": 1},
					"surfshark": {"version": 1, "timestamp": 1},
					"torguard": {"version": 1, "timestamp": 1},
					"vpn unlimited": {"version": 1, "timestamp": 1},
					"vyprvpn": {"version": 1, "timestamp": 1},
					"wevpn": {"version": 1, "timestamp": 1},
					"windscribe": {"version": 1, "timestamp": 1}
				}`),
			hardcoded: populateProviders(1, 0, models.AllServers{}),
			persisted: populateProviders(1, 1, models.AllServers{}),
		},
		"different versions": {
			b: []byte(`{
				"cyberghost": {"version": 1, "timestamp": 1},
				"expressvpn": {"version": 1, "timestamp": 1},
				"fastestvpn": {"version": 1, "timestamp": 1},
				"hidemyass": {"version": 1, "timestamp": 1},
				"ipvanish": {"version": 1, "timestamp": 1},
				"ivpn": {"version": 1, "timestamp": 1},
				"mullvad": {"version": 1, "timestamp": 1},
				"nordvpn": {"version": 1, "timestamp": 1},
				"perfect privacy": {"version": 1, "timestamp": 1},
				"privado": {"version": 1, "timestamp": 1},
				"private internet access": {"version": 1, "timestamp": 1},
				"privatevpn": {"version": 1, "timestamp": 1},
				"protonvpn": {"version": 1, "timestamp": 1},
				"purevpn": {"version": 1, "timestamp": 1},
				"surfshark": {"version": 1, "timestamp": 1},
				"torguard": {"version": 1, "timestamp": 1},
				"vpn unlimited": {"version": 1, "timestamp": 1},
				"vyprvpn": {"version": 1, "timestamp": 1},
				"wevpn": {"version": 1, "timestamp": 1},
				"windscribe": {"version": 1, "timestamp": 1}
			}`),
			hardcoded: populateProviders(2, 0, models.AllServers{}),
			logged: []string{
				"Cyberghost servers from file discarded because they have version 1 and hardcoded servers have version 2",
				"Expressvpn servers from file discarded because they have version 1 and hardcoded servers have version 2",
				"Fastestvpn servers from file discarded because they have version 1 and hardcoded servers have version 2",
				"Hidemyass servers from file discarded because they have version 1 and hardcoded servers have version 2",
				"Ipvanish servers from file discarded because they have version 1 and hardcoded servers have version 2",
				"Ivpn servers from file discarded because they have version 1 and hardcoded servers have version 2",
				"Mullvad servers from file discarded because they have version 1 and hardcoded servers have version 2",
				"Nordvpn servers from file discarded because they have version 1 and hardcoded servers have version 2",
				"Perfect Privacy servers from file discarded because they have version 1 and hardcoded servers have version 2",
				"Privado servers from file discarded because they have version 1 and hardcoded servers have version 2",
				"Private Internet Access servers from file discarded because they have version 1 and hardcoded servers have version 2", //nolint:lll
				"Privatevpn servers from file discarded because they have version 1 and hardcoded servers have version 2",
				"Protonvpn servers from file discarded because they have version 1 and hardcoded servers have version 2",
				"Purevpn servers from file discarded because they have version 1 and hardcoded servers have version 2",
				"Surfshark servers from file discarded because they have version 1 and hardcoded servers have version 2",
				"Torguard servers from file discarded because they have version 1 and hardcoded servers have version 2",
				"Vpn Unlimited servers from file discarded because they have version 1 and hardcoded servers have version 2",
				"Vyprvpn servers from file discarded because they have version 1 and hardcoded servers have version 2",
				"Wevpn servers from file discarded because they have version 1 and hardcoded servers have version 2",
				"Windscribe servers from file discarded because they have version 1 and hardcoded servers have version 2",
			},
			persisted: models.AllServers{
				ProviderToServers: map[string]models.Servers{},
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			logger := NewMockInfoErrorer(ctrl)
			var previousLogCall *gomock.Call
			for _, logged := range testCase.logged {
				call := logger.EXPECT().Info(logged)
				if previousLogCall != nil {
					call.After(previousLogCall)
				}
				previousLogCall = call
			}

			s := &Storage{
				logger: logger,
			}

			servers, err := s.extractServersFromBytes(testCase.b, testCase.hardcoded)

			if testCase.errMessage != "" {
				assert.EqualError(t, err, testCase.errMessage)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.persisted, servers)
		})
	}

	t.Run("hardcoded panic", func(t *testing.T) {
		t.Parallel()

		s := &Storage{}

		allProviders := providers.All()
		require.GreaterOrEqual(t, len(allProviders), 2)

		b := []byte(`{}`)
		hardcoded := models.AllServers{
			ProviderToServers: map[string]models.Servers{
				allProviders[0]: {},
				// Missing provider allProviders[1]
			},
		}
		expectedPanicValue := fmt.Sprintf("provider %s not found in hardcoded servers map", allProviders[1])
		assert.PanicsWithValue(t, expectedPanicValue, func() {
			_, _ = s.extractServersFromBytes(b, hardcoded)
		})
	})
}
