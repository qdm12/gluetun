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

	type stringCallWithWarning struct {
		call    bool
		s       string
		warning string
		err     error
	}

	type durationCall struct {
		call     bool
		duration time.Duration
		err      error
	}

	testCases := map[string]struct {
		serverAddress stringCallWithWarning
		vpnInitial    durationCall
		vpnAddition   durationCall
		expected      Health
		err           error
	}{
		"success": {
			serverAddress: stringCallWithWarning{
				call: true,
				s:    "127.0.0.1:9999",
			},
			vpnInitial: durationCall{
				call:     true,
				duration: time.Second,
			},
			vpnAddition: durationCall{
				call:     true,
				duration: time.Minute,
			},
			expected: Health{
				ServerAddress: "127.0.0.1:9999",
				VPN: HealthyWait{
					Initial:  time.Second,
					Addition: time.Minute,
				},
			},
		},
		"listening address error": {
			serverAddress: stringCallWithWarning{
				call:    true,
				s:       "127.0.0.1:9999",
				warning: "warning",
				err:     errDummy,
			},
			expected: Health{
				ServerAddress: "127.0.0.1:9999",
			},
			err: errors.New("environment variable HEALTH_SERVER_ADDRESS: dummy"),
		},
		"initial error": {
			serverAddress: stringCallWithWarning{
				call: true,
			},
			vpnInitial: durationCall{
				call:     true,
				duration: time.Second,
				err:      errDummy,
			},
			expected: Health{
				VPN: HealthyWait{
					Initial: time.Second,
				},
			},
			err: errors.New("environment variable HEALTH_VPN_DURATION_INITIAL: dummy"),
		},
		"addition error": {
			serverAddress: stringCallWithWarning{
				call: true,
			},
			vpnInitial: durationCall{
				call:     true,
				duration: time.Second,
			},
			vpnAddition: durationCall{
				call:     true,
				duration: time.Minute,
				err:      errDummy,
			},
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

			if testCase.serverAddress.call {
				value := testCase.serverAddress.s
				warning := testCase.serverAddress.warning
				err := testCase.serverAddress.err
				env.EXPECT().ListeningAddress("HEALTH_SERVER_ADDRESS", gomock.Any()).
					Return(value, warning, err)
				if warning != "" {
					logger.EXPECT().Warn("environment variable HEALTH_SERVER_ADDRESS: " + warning)
				}
			}

			if testCase.vpnInitial.call {
				value := testCase.vpnInitial.duration
				err := testCase.vpnInitial.err
				env.EXPECT().
					Duration("HEALTH_VPN_DURATION_INITIAL", gomock.Any()).
					Return(value, err)
			}

			if testCase.vpnAddition.call {
				value := testCase.vpnAddition.duration
				err := testCase.vpnAddition.err
				env.EXPECT().
					Duration("HEALTH_VPN_DURATION_ADDITION", gomock.Any()).
					Return(value, err)
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
