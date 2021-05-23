package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CipherLines(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		version string
		lines   []string
	}{
		"empty version": {
			lines: []string{
				"data-ciphers-fallback AES",
				"data-ciphers AES",
			},
		},
		"2.4.5": {
			version: "2.4.5",
			lines:   []string{"cipher AES"},
		},
		"2.5.3": {
			version: "2.5.3",
			lines: []string{
				"data-ciphers-fallback AES",
				"data-ciphers AES",
			},
		},
	}
	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			const cipher = "AES"

			lines := CipherLines(cipher, testCase.version)

			assert.Equal(t, testCase.lines, lines)
		})
	}
}
