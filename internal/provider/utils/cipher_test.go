package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CipherLines(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		ciphers []string
		version string
		lines   []string
	}{
		"empty version": {
			ciphers: []string{"AES"},
			lines: []string{
				"data-ciphers-fallback AES",
				"data-ciphers AES",
			},
		},
		"2.5": {
			ciphers: []string{"AES", "CBC"},
			version: "2.5",
			lines: []string{
				"data-ciphers-fallback AES",
				"data-ciphers AES:CBC",
			},
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			lines := CipherLines(testCase.ciphers)

			assert.Equal(t, testCase.lines, lines)
		})
	}
}
