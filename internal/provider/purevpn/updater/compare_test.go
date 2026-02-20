package updater

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_comparePlaceNames(t *testing.T) {
	t.Parallel() // Allow the top-level test to run in parallel

	testCases := map[string]struct {
		a    string
		b    string
		want bool
	}{
		"exact_match": {
			a:    "Paris",
			b:    "Paris",
			want: true,
		},
		"difference_in_casing_and_whitespace": {
			a:    "  Montreal",
			b:    "montreal  ",
			want: true,
		},
		"accent_normalization": {
			a:    "Montréal",
			b:    "Montreal",
			want: true,
		},
		"single_character_typo": {
			a:    "Lyon",
			b:    "Lyonn",
			want: true,
		},
		"too_many_differences": {
			a:    "London",
			b:    "Londres",
			want: false,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			result := comparePlaceNames(testCase.a, testCase.b)
			assert.Equal(t, testCase.want, result)
		})
	}
}
