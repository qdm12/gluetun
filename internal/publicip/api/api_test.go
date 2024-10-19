package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParseProvider(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		s          string
		provider   Provider
		errWrapped error
		errMessage string
	}{
		"empty": {
			errWrapped: ErrProviderNotValid,
			errMessage: `API name is not valid: "" can only be "cloudflare", "ifconfigco", "ip2location" or "ipinfo"`,
		},
		"invalid": {
			s:          "xyz",
			errWrapped: ErrProviderNotValid,
			errMessage: `API name is not valid: "xyz" can only be "cloudflare", "ifconfigco", "ip2location" or "ipinfo"`,
		},
		"ipinfo": {
			s:        "ipinfo",
			provider: IPInfo,
		},
		"IpInfo": {
			s:        "IpInfo",
			provider: IPInfo,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			provider, err := ParseProvider(testCase.s)

			assert.Equal(t, testCase.provider, provider)
			assert.ErrorIs(t, err, testCase.errWrapped)
			if testCase.errWrapped != nil {
				assert.EqualError(t, err, testCase.errMessage)
			}
		})
	}
}
