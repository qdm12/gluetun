package configuration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Provider_vpnUnlimitedLines(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		settings Provider
		lines    []string
	}{
		"empty settings": {},
		"full settings": {
			settings: Provider{
				ServerSelection: ServerSelection{
					Countries: []string{"A", "B"},
					Cities:    []string{"C", "D"},
					Hostnames: []string{"E", "F"},
				},
			},
			lines: []string{
				"|--Countries: A, B",
				"|--Cities: C, D",
				"|--Hostnames: E, F",
			},
		},
	}
	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			lines := testCase.settings.vpnUnlimitedLines()

			assert.Equal(t, testCase.lines, lines)
		})
	}
}
