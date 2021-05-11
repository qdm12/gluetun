package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FilterByPossibilities(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		value         string
		possibilities []string
		filtered      bool
	}{
		"no possibilities": {},
		"value not in possibilities": {
			value:         "c",
			possibilities: []string{"a", "b"},
			filtered:      true,
		},
		"value in possibilities": {
			value:         "c",
			possibilities: []string{"a", "b", "c"},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			filtered := FilterByPossibilities(testCase.value, testCase.possibilities)
			assert.Equal(t, testCase.filtered, filtered)
		})
	}
}
