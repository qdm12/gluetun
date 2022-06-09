package updater

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func Test_parseFilename(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		fileName string
		hostname string
		country  string
		city     string
		err      error
	}{
		"unknown country code": {
			fileName: "ipvanish-unknown-host.ovpn",
			hostname: "host.ipvanish.com",
			err:      errors.New("country code is unknown: unknown"),
		},
		"country code only": {
			fileName: "ipvanish-ca-host.ovpn",
			hostname: "host.ipvanish.com",
			country:  "Canada",
		},
		"country code and city": {
			fileName: "ipvanish-ca-sao-paulo-host.ovpn",
			hostname: "host.ipvanish.com",
			country:  "Canada",
			city:     "Sao Paulo",
		},
	}
	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			titleCaser := cases.Title(language.English)
			country, city, err := parseFilename(testCase.fileName, testCase.hostname, titleCaser)

			if testCase.err != nil {
				require.Error(t, err)
				assert.Equal(t, testCase.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, testCase.country, country)
			assert.Equal(t, testCase.city, city)
		})
	}
}
