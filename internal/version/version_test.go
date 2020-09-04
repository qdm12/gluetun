package version

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_formatDuration(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		duration time.Duration
		s        string
	}{
		"zero": {
			s: "0 second",
		},
		"one second": {
			duration: time.Second,
			s:        "1 second",
		},
		"59 seconds": {
			duration: 59 * time.Second,
			s:        "59 seconds",
		},
		"1 minute": {
			duration: time.Minute,
			s:        "1 minute",
		},
		"2 minutes": {
			duration: 2 * time.Minute,
			s:        "2 minutes",
		},
		"1 hour": {
			duration: time.Hour,
			s:        "60 minutes",
		},
		"2 hours": {
			duration: 2 * time.Hour,
			s:        "2 hours",
		},
		"26 hours": {
			duration: 26 * time.Hour,
			s:        "26 hours",
		},
		"28 hours": {
			duration: 28 * time.Hour,
			s:        "28 hours",
		},
		"55 hours": {
			duration: 55 * time.Hour,
			s:        "2 days",
		},
	}
	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			s := formatDuration(testCase.duration)
			assert.Equal(t, testCase.s, s)
		})
	}
}
