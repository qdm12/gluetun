package configuration

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_HealthyWait_String(t *testing.T) {
	t.Parallel()

	var healthyWait HealthyWait
	const expected = "|--Initial duration: 0s"

	s := healthyWait.String()

	assert.Equal(t, expected, s)
}

func Test_HealthyWait_lines(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		settings HealthyWait
		lines    []string
	}{
		"empty": {
			lines: []string{
				"|--Initial duration: 0s",
			},
		},
		"filled settings": {
			settings: HealthyWait{
				Initial:  time.Second,
				Addition: time.Minute,
			},
			lines: []string{
				"|--Initial duration: 1s",
				"|--Addition duration: 1m0s",
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			lines := testCase.settings.lines()

			assert.Equal(t, testCase.lines, lines)
		})
	}
}
