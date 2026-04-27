package storage

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func populateProviderToVersion(providerToVersion map[string]uint16) map[string]uint16 {
	allProviders := providers.All()
	for _, provider := range allProviders {
		_, has := providerToVersion[provider]
		if has {
			continue
		}

		providerToVersion[provider] = 0
	}
	return providerToVersion
}

func Test_extractServersFromBytes(t *testing.T) {
	t.Parallel()

	type logLine struct {
		level   log.Level
		message string
	}

	testCases := map[string]struct {
		b                 []byte
		hardcodedVersions map[string]uint16
		logged            []logLine
		persisted         models.AllServers
		errMessage        string
	}{
		"bad JSON": {
			b:          []byte("garbage"),
			errMessage: "decoding servers: invalid character 'g' looking for beginning of value",
		},
		"bad provider JSON": {
			b:                 []byte(`{"cyberghost": "garbage"}`),
			hardcodedVersions: populateProviderToVersion(map[string]uint16{}),
			errMessage: "decoding servers version for provider Cyberghost: " +
				"json: cannot unmarshal string into Go value of type struct { Version uint16 \"json:\\\"version\\\"\"; " +
				"Timestamp int64 \"json:\\\"timestamp\\\"\"; " +
				"Filepath string \"json:\\\"filepath\\\"\" }",
		},
		"bad servers array JSON": {
			b: []byte(`{"cyberghost": {"version": 1, "servers": "garbage"}}`),
			hardcodedVersions: populateProviderToVersion(map[string]uint16{
				providers.Cyberghost: 1,
			}),
			errMessage: "decoding servers for provider Cyberghost: " +
				"json: cannot unmarshal string into Go struct field Servers.servers of type []models.Server",
		},
		"absent provider keys": {
			b: []byte(`{}`),
			hardcodedVersions: populateProviderToVersion(map[string]uint16{
				providers.Cyberghost: 1,
			}),
			persisted: models.AllServers{
				ProviderToServers: map[string]models.Servers{},
			},
		},
		"same versions": {
			b: []byte(`{
					"cyberghost": {"version": 1, "timestamp": 0}
				}`),
			hardcodedVersions: populateProviderToVersion(map[string]uint16{
				providers.Cyberghost: 1,
			}),
			persisted: models.AllServers{
				ProviderToServers: map[string]models.Servers{
					providers.Cyberghost: {Version: 1},
				},
			},
		},
		"different versions": {
			b: []byte(`{
				"cyberghost": {"version": 1, "timestamp": 1}
			}`),
			hardcodedVersions: populateProviderToVersion(map[string]uint16{
				providers.Cyberghost: 2,
			}),
			logged: []logLine{
				{level: log.LevelInfo, message: "Cyberghost servers from file discarded because they have version 1" +
					" and hardcoded servers have version 2"},
			},
			persisted: models.AllServers{
				ProviderToServers: map[string]models.Servers{},
			},
		},
		"preferred_different_versions": {
			b: []byte(`{
				"cyberghost": {"version": 1, "timestamp": 1, "preferred": true}
			}`),
			hardcodedVersions: populateProviderToVersion(map[string]uint16{
				providers.Cyberghost: 2,
			}),
			logged: []logLine{
				{level: log.LevelWarn, message: "Cyberghost preferred servers from file discarded because they have version 1" +
					" and hardcoded servers have version 2"},
			},
			persisted: models.AllServers{
				ProviderToServers: map[string]models.Servers{},
			},
			errMessage: "",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			logger := NewMockLogger(ctrl)
			var previousLogCall *gomock.Call
			for _, logged := range testCase.logged {
				var call *gomock.Call
				switch logged.level { //nolint:exhaustive
				case log.LevelInfo:
					call = logger.EXPECT().Info(logged.message)
				case log.LevelWarn:
					call = logger.EXPECT().Warn(logged.message)
				default:
					t.Fatalf("invalid log level %d in test case", logged.level)
				}
				if previousLogCall != nil {
					call.After(previousLogCall)
				}
				previousLogCall = call
			}

			s := &Storage{
				logger: logger,
			}

			servers, err := s.extractServersFromBytes(testCase.b, testCase.hardcodedVersions)

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
		hardcodedVersions := map[string]uint16{
			allProviders[0]: 1,
			// Missing provider allProviders[1]
		}
		expectedPanicValue := fmt.Sprintf("provider %s not found in hardcoded servers map; "+
			"did you add the provider key in the embedded servers.json?", allProviders[1])
		assert.PanicsWithValue(t, expectedPanicValue, func() {
			_, _ = s.extractServersFromBytes(b, hardcodedVersions)
		})
	})
}
