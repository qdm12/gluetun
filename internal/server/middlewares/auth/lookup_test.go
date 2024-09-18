package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Read reads the toml file specified by the filepath given.
func Test_settingsToLookupMap(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		settings     Settings
		routeToRoles map[string][]internalRole
		errWrapped   error
		errMessage   string
	}{
		"empty_settings": {
			routeToRoles: map[string][]internalRole{},
		},
		"auth_method_not_supported": {
			settings: Settings{
				Roles: []Role{{Name: "a", Auth: "bad"}},
			},
			errWrapped: ErrMethodNotSupported,
			errMessage: "authentication method not supported: bad",
		},
		"success": {
			settings: Settings{
				Roles: []Role{
					{Name: "a", Auth: AuthNone, Routes: []string{"GET /path"}},
					{Name: "b", Auth: AuthNone, Routes: []string{"GET /path", "PUT /path"}},
				},
			},
			routeToRoles: map[string][]internalRole{
				"GET /path": {
					{name: "a", checker: newNoneMethod()}, // deduplicated method
				},
				"PUT /path": {
					{name: "b", checker: newNoneMethod()},
				}},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			routeToRoles, err := settingsToLookupMap(testCase.settings)

			assert.Equal(t, testCase.routeToRoles, routeToRoles)
			assert.ErrorIs(t, err, testCase.errWrapped)
			if testCase.errWrapped != nil {
				assert.EqualError(t, err, testCase.errMessage)
			}
		})
	}
}
