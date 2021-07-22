package configuration

import (
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/qdm12/golibs/params/mock_params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Health_String(t *testing.T) {
	t.Parallel()

	var health Health
	const expected = "|--Health:\n   |--OpenVPN:\n      |--Initial duration: 0s"

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
				"   |--OpenVPN:",
				"      |--Initial duration: 0s",
			},
		},
		"filled settings": {
			settings: Health{
				OpenVPN: HealthyWait{
					Initial:  time.Second,
					Addition: time.Minute,
				},
			},
			lines: []string{
				"|--Health:",
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
		expected                Health
		err                     error
	}{
		"success": {
			openvpnInitialDuration:  time.Second,
			openvpnAdditionDuration: time.Minute,
			expected: Health{
				OpenVPN: HealthyWait{
					Initial:  time.Second,
					Addition: time.Minute,
				},
			},
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
			env.EXPECT().
				Duration("HEALTH_OPENVPN_DURATION_INITIAL", gomock.Any()).
				Return(testCase.openvpnInitialDuration, testCase.openvpnInitialErr)
			if testCase.openvpnInitialErr == nil {
				env.EXPECT().
					Duration("HEALTH_OPENVPN_DURATION_ADDITION", gomock.Any()).
					Return(testCase.openvpnAdditionDuration, testCase.openvpnAdditionErr)
			}

			r := reader{
				env: env,
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
