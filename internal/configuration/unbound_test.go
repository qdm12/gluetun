package configuration

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/qdm12/dns/pkg/provider"
	"github.com/qdm12/dns/pkg/unbound"
	"github.com/qdm12/golibs/params/mock_params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_DNS_readUnboundProviders(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		envValue string
		envErr   error
		expected DNS
		err      error
	}{
		"bad value": {
			envValue: "invalid",
			err:      errors.New(`invalid DNS over TLS provider: cannot parse provider: "invalid"`),
		},
		"env error": {
			envErr: errors.New("env error"),
			err:    errors.New("environment variable DOT_PROVIDERS: env error"),
		},
		"multiple valid values": {
			envValue: "cloudflare,google",
			expected: DNS{
				Unbound: unbound.Settings{
					Providers: []provider.Provider{
						provider.Cloudflare(),
						provider.Google(),
					},
				},
			},
		},
		"one invalid value in two": {
			envValue: "cloudflare,invalid",
			expected: DNS{
				Unbound: unbound.Settings{
					Providers: []provider.Provider{
						provider.Cloudflare(),
					},
				},
			},
			err: errors.New(`invalid DNS over TLS provider: cannot parse provider: "invalid"`),
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			env := mock_params.NewMockInterface(ctrl)
			env.EXPECT().Get("DOT_PROVIDERS", gomock.Any()).
				Return(testCase.envValue, testCase.envErr)

			var settings DNS
			err := settings.readUnboundProviders(env)

			if testCase.err != nil {
				require.Error(t, err)
				assert.Equal(t, testCase.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.expected, settings)
		})
	}
}
