package ivpn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseFilename(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		fileName string
		country  string
		city     string
	}{
		"empty filename": {},
		"country only": {
			fileName: "Country.ovpn",
			country:  "Country",
		},
		"country and city": {
			fileName: "Country-City.ovpn",
			country:  "Country",
			city:     "City",
		},
		"composite country and city": {
			fileName: "Coun_try-Ci_ty.ovpn",
			country:  "Coun try",
			city:     "Ci ty",
		},
	}
	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			country, city := parseFilename(testCase.fileName)
			assert.Equal(t, testCase.country, country)
			assert.Equal(t, testCase.city, city)
		})
	}
}
