package updater

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_comparePlaceNames(t *testing.T) {
	t.Parallel() // Allow the top-level test to run in parallel

	tests := map[string]struct {
		inputA   string
		inputB   string
		expected bool
	}{
		"exact_match": {
			inputA:   "Paris",
			inputB:   "Paris",
			expected: true,
		},
		"difference_in_casing_and_whitespace": {
			inputA:   "  Montreal",
			inputB:   "montreal  ",
			expected: true,
		},
		"accent_normalization": {
			inputA:   "Montréal",
			inputB:   "Montreal",
			expected: true,
		},
		"single_character_typo": {
			inputA:   "Lyon",
			inputB:   "Lyonn",
			expected: true,
		},
		"too_many_differences": {
			inputA:   "London",
			inputB:   "Londres",
			expected: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			result := comparePlaceNames(tc.inputA, tc.inputB)
			assert.Equal(t, tc.expected, result)
		})
	}
}
