package configuration

import (
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/qdm12/golibs/logging/mock_logging"
	"github.com/qdm12/golibs/params/mock_params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Health_String(t *testing.T) {
	t.Parallel()

	var health Health
	const expected = "|--Health:\n   |--Server address: \n   |--OpenVPN:\n      |--Initial duration: 0s"

	s := health.String()

	assert.Equal(t, expected, s)
}

func Test_Health_lines(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		settings Health
		lines    []string
	}{
		"empty": {
			lines: []string{
				"|--Health:",
				"   |--Server address: ",
				"   |--OpenVPN:",
				"      |--Initial duration: 0s",
			},
		},
		"filled settings": {
			settings: Health{
				ServerAddress: "address:9999",
				OpenVPN: HealthyWait{
					Initial:  time.Second,
					Addition: time.Minute,
				},
			},
			lines: []string{
				"|--Health:",
				"   |--Server address: address:9999",
				"   |--OpenVPN:",
				"      |--Initial duration: 1s",
				"      |--Addition duration: 1m0s",
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

func Test_Health_read(t *testing.T) {
	t.Parallel()

	errDummy := errors.New("dummy")

	testCases := map[string]struct {
		openvpnInitialDuration  time.Duration
		openvpnInitialErr       error
		openvpnAdditionDuration time.Duration
		openvpnAdditionErr      error
		serverAddress           string
		serverAddressWarning    string
		serverAddressErr        error
		expected                Health
		err                     error
	}{
		"success": {
			openvpnInitialDuration:  time.Second,
			openvpnAdditionDuration: time.Minute,
			serverAddress:           "127.0.0.1:9999",
			expected: Health{
				ServerAddress: "127.0.0.1:9999",
				OpenVPN: HealthyWait{
					Initial:  time.Second,
					Addition: time.Minute,
				},
			},
		},
		"listening address error": {
			openvpnInitialDuration:  time.Second,
			openvpnAdditionDuration: time.Minute,
			serverAddress:           "127.0.0.1:9999",
			serverAddressWarning:    "warning",
			serverAddressErr:        errDummy,
			expected: Health{
				ServerAddress: "127.0.0.1:9999",
			},
			err: errDummy,
		},
		"initial error": {
			openvpnInitialDuration:  time.Second,
			openvpnInitialErr:       errDummy,
			openvpnAdditionDuration: time.Minute,
			expected: Health{
				OpenVPN: HealthyWait{
					Initial: time.Second,
				},
			},
			err: errDummy,
		},
		"addition error": {
			openvpnInitialDuration:  time.Second,
			openvpnAdditionDuration: time.Minute,
			openvpnAdditionErr:      errDummy,
			expected: Health{
				OpenVPN: HealthyWait{
					Initial:  time.Second,
					Addition: time.Minute,
				},
			},
			err: errDummy,
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			env := mock_params.NewMockEnv(ctrl)
			logger := mock_logging.NewMockLogger(ctrl)

			env.EXPECT().ListeningAddress("HEALTH_SERVER_ADDRESS", gomock.Any()).
				Return(testCase.serverAddress, testCase.serverAddressWarning,
					testCase.serverAddressErr)
			if testCase.serverAddressWarning != "" {
				logger.EXPECT().Warn("health server address: " + testCase.serverAddressWarning)
			}

			if testCase.serverAddressErr == nil {
				env.EXPECT().
					Duration("HEALTH_OPENVPN_DURATION_INITIAL", gomock.Any()).
					Return(testCase.openvpnInitialDuration, testCase.openvpnInitialErr)
				if testCase.openvpnInitialErr == nil {
					env.EXPECT().
						Duration("HEALTH_OPENVPN_DURATION_ADDITION", gomock.Any()).
						Return(testCase.openvpnAdditionDuration, testCase.openvpnAdditionErr)
				}
			}

			r := reader{
				env:    env,
				logger: logger,
			}

			var health Health

			err := health.read(r)

			if testCase.err != nil {
				require.Error(t, err)
				assert.Equal(t, testCase.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.expected, health)
		})
	}
}
