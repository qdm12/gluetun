package format

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_FriendlyDuration(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		duration time.Duration
		friendly string
	}{
		"zero": {
			friendly: "0 second",
		},
		"one_second": {
			duration: time.Second,
			friendly: "1 second",
		},
		"59_seconds": {
			duration: 59 * time.Second,
			friendly: "59 seconds",
		},
		"1_minute": {
			duration: time.Minute,
			friendly: "1 minute",
		},
		"2_minutes": {
			duration: 2 * time.Minute,
			friendly: "2 minutes",
		},
		"1_hour": {
			duration: time.Hour,
			friendly: "60 minutes",
		},
		"2_hours": {
			duration: 2 * time.Hour,
			friendly: "2 hours",
		},
		"26_hours": {
			duration: 26 * time.Hour,
			friendly: "26 hours",
		},
		"28_hours": {
			duration: 28 * time.Hour,
			friendly: "28 hours",
		},
		"55_hours": {
			duration: 55 * time.Hour,
			friendly: "2 days",
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := FriendlyDuration(testCase.duration)

			assert.Equal(t, testCase.friendly, s)
		})
	}
}
