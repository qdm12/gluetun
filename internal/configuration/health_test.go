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
	const expected = "|--Health:\n   |--Server address: \n   |--VPN:\n      |--Initial duration: 0s"

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
				"   |--VPN:",
				"      |--Initial duration: 0s",
			},
		},
		"filled settings": {
			settings: Health{
				ServerAddress: "address:9999",
				VPN: HealthyWait{
					Initial:  time.Second,
					Addition: time.Minute,
				},
			},
			lines: []string{
				"|--Health:",
				"   |--Server address: address:9999",
				"   |--VPN:",
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
		vpnInitialDuration   time.Duration
		vpnInitialErr        error
		vpnAdditionDuration  time.Duration
		vpnAdditionErr       error
		serverAddress        string
		serverAddressWarning string
		serverAddressErr     error
		expected             Health
		err                  error
	}{
		"success": {
			vpnInitialDuration:  time.Second,
			vpnAdditionDuration: time.Minute,
			serverAddress:       "127.0.0.1:9999",
			expected: Health{
				ServerAddress: "127.0.0.1:9999",
				VPN: HealthyWait{
					Initial:  time.Second,
					Addition: time.Minute,
				},
			},
		},
		"listening address error": {
			vpnInitialDuration:   time.Second,
			vpnAdditionDuration:  time.Minute,
			serverAddress:        "127.0.0.1:9999",
			serverAddressWarning: "warning",
			serverAddressErr:     errDummy,
			expected: Health{
				ServerAddress: "127.0.0.1:9999",
			},
			err: errors.New("environment variable HEALTH_SERVER_ADDRESS: dummy"),
		},
		"initial error": {
			vpnInitialDuration:  time.Second,
			vpnInitialErr:       errDummy,
			vpnAdditionDuration: time.Minute,
			expected: Health{
				VPN: HealthyWait{
					Initial: time.Second,
				},
			},
			err: errors.New("environment variable HEALTH_VPN_DURATION_INITIAL: dummy"),
		},
		"addition error": {
			vpnInitialDuration:  time.Second,
			vpnAdditionDuration: time.Minute,
			vpnAdditionErr:      errDummy,
			expected: Health{
				VPN: HealthyWait{
					Initial:  time.Second,
					Addition: time.Minute,
				},
			},
			err: errors.New("environment variable HEALTH_VPN_DURATION_ADDITION: dummy"),
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			env := mock_params.NewMockInterface(ctrl)
			logger := mock_logging.NewMockLogger(ctrl)

			env.EXPECT().ListeningAddress("HEALTH_SERVER_ADDRESS", gomock.Any()).
				Return(testCase.serverAddress, testCase.serverAddressWarning,
					testCase.serverAddressErr)
			if testCase.serverAddressWarning != "" {
				logger.EXPECT().Warn("environment variable HEALTH_SERVER_ADDRESS: " + testCase.serverAddressWarning)
			}

			if testCase.serverAddressErr == nil {
				env.EXPECT().
					Duration("HEALTH_VPN_DURATION_INITIAL", gomock.Any()).
					Return(testCase.vpnInitialDuration, testCase.vpnInitialErr)
				if testCase.vpnInitialErr == nil {
					env.EXPECT().
						Duration("HEALTH_VPN_DURATION_ADDITION", gomock.Any()).
						Return(testCase.vpnAdditionDuration, testCase.vpnAdditionErr)
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
