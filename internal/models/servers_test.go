package models

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_AllServers_MarshalJSON(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		allServers *AllServers
		dataString string
		errWrapped error
		errMessage string
	}{
		"empty": {
			allServers: &AllServers{
				ProviderToServers: map[string]Servers{},
			},
			dataString: `{"version":0,` +
				`"cyberghost":{"version":0,"timestamp":0},` +
				`"expressvpn":{"version":0,"timestamp":0},` +
				`"fastestvpn":{"version":0,"timestamp":0},` +
				`"hidemyass":{"version":0,"timestamp":0},` +
				`"ipvanish":{"version":0,"timestamp":0},` +
				`"ivpn":{"version":0,"timestamp":0},` +
				`"mullvad":{"version":0,"timestamp":0},` +
				`"nordvpn":{"version":0,"timestamp":0},` +
				`"perfect privacy":{"version":0,"timestamp":0},` +
				`"privado":{"version":0,"timestamp":0},` +
				`"private internet access":{"version":0,"timestamp":0},` +
				`"privatevpn":{"version":0,"timestamp":0},` +
				`"protonvpn":{"version":0,"timestamp":0},` +
				`"purevpn":{"version":0,"timestamp":0},` +
				`"surfshark":{"version":0,"timestamp":0},` +
				`"torguard":{"version":0,"timestamp":0},` +
				`"vpn unlimited":{"version":0,"timestamp":0},` +
				`"vyprvpn":{"version":0,"timestamp":0},` +
				`"wevpn":{"version":0,"timestamp":0},` +
				`"windscribe":{"version":0,"timestamp":0}}`,
		},
		"two known providers": {
			allServers: &AllServers{
				Version: 1,
				ProviderToServers: map[string]Servers{
					providers.Cyberghost: {
						Version:   1,
						Timestamp: 1000,
						Servers: []Server{
							{Country: "A"},
							{Country: "B"},
						},
					},
					providers.Privado: {
						Version:   2,
						Timestamp: 2000,
						Servers: []Server{
							{City: "C"},
							{City: "D"},
						},
					},
				},
			},
			dataString: `{"version":1,` +
				`"cyberghost":{"version":1,"timestamp":1000,"servers":[{"country":"A"},{"country":"B"}]},` +
				`"expressvpn":{"version":0,"timestamp":0},` +
				`"fastestvpn":{"version":0,"timestamp":0},` +
				`"hidemyass":{"version":0,"timestamp":0},` +
				`"ipvanish":{"version":0,"timestamp":0},` +
				`"ivpn":{"version":0,"timestamp":0},` +
				`"mullvad":{"version":0,"timestamp":0},` +
				`"nordvpn":{"version":0,"timestamp":0},` +
				`"perfect privacy":{"version":0,"timestamp":0},` +
				`"privado":{"version":2,"timestamp":2000,"servers":[{"city":"C"},{"city":"D"}]},` +
				`"private internet access":{"version":0,"timestamp":0},` +
				`"privatevpn":{"version":0,"timestamp":0},` +
				`"protonvpn":{"version":0,"timestamp":0},` +
				`"purevpn":{"version":0,"timestamp":0},` +
				`"surfshark":{"version":0,"timestamp":0},` +
				`"torguard":{"version":0,"timestamp":0},` +
				`"vpn unlimited":{"version":0,"timestamp":0},` +
				`"vyprvpn":{"version":0,"timestamp":0},` +
				`"wevpn":{"version":0,"timestamp":0},` +
				`"windscribe":{"version":0,"timestamp":0}}`,
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Populate all providers in all servers
			for _, provider := range providers.All() {
				_, has := testCase.allServers.ProviderToServers[provider]
				if !has {
					testCase.allServers.ProviderToServers[provider] = Servers{}
				}
			}

			data, err := testCase.allServers.MarshalJSON()
			assert.ErrorIs(t, err, testCase.errWrapped)
			if err != nil {
				assert.EqualError(t, err, testCase.errMessage)
			}
			assert.Equal(t, testCase.dataString, string(data))

			data, err = json.Marshal(testCase.allServers)
			assert.ErrorIs(t, err, testCase.errWrapped)
			if err != nil {
				assert.EqualError(t, err, testCase.errMessage)
			}
			assert.Equal(t, testCase.dataString, string(data))

			buffer := bytes.NewBuffer(nil)
			encoder := json.NewEncoder(buffer)
			// encoder.SetIndent("", "  ")
			err = encoder.Encode(testCase.allServers)
			require.NoError(t, err)
			assert.Equal(t, testCase.dataString+"\n", buffer.String())
		})
	}
}

func Test_AllServers_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		dataString string
		allServers AllServers
		errWrapped error
		errMessage string
	}{
		"empty": {
			dataString: "{}",
			allServers: AllServers{},
		},
		"two known providers": {
			dataString: `{"version":1,` +
				`"cyberghost":{"version":1,"timestamp":1000,"servers":[{"country":"A"},{"country":"B"}]},` +
				`"privado":{"version":2,"timestamp":2000,"servers":[{"city":"C"},{"city":"D"}]}}`,
			allServers: AllServers{
				Version: 1,
				ProviderToServers: map[string]Servers{
					providers.Cyberghost: {
						Version:   1,
						Timestamp: 1000,
						Servers: []Server{
							{Country: "A"},
							{Country: "B"},
						},
					},
					providers.Privado: {
						Version:   2,
						Timestamp: 2000,
						Servers: []Server{
							{City: "C"},
							{City: "D"},
						},
					},
				},
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			data := []byte(testCase.dataString)
			var allServers AllServers

			err := json.Unmarshal(data, &allServers)

			assert.ErrorIs(t, err, testCase.errWrapped)
			if err != nil {
				assert.EqualError(t, err, testCase.errMessage)
			}
			assert.Equal(t, testCase.allServers, allServers)
		})
	}
}
