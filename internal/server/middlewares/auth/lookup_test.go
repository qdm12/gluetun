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
				Auths: []Auth{{Name: "bad", Method: "not_supported"}},
			},
			errWrapped: ErrMethodNotSupported,
			errMessage: "authentication method not supported: bad",
		},
		"auth_name_not_defined": {
			settings: Settings{
				Auths: []Auth{{Name: "x", Method: MethodNone}, {Name: "y", Method: MethodNone}},
				Roles: []Role{{Name: "a", Auths: []string{"xyz"}}},
			},
			errWrapped: ErrAuthNameNotDefined,
			errMessage: "authentication name not defined: xyz is not one of x or y",
		},
		"success": {
			settings: Settings{
				Auths: []Auth{
					{Name: "x", Method: MethodNone},
					{Name: "y", Method: MethodNone},
				},
				Roles: []Role{
					{Name: "a", Auths: []string{"x"}, Routes: []string{"GET /path"}},
					{Name: "b", Auths: []string{"x", "y"}, Routes: []string{"GET /path", "PUT /path"}},
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
