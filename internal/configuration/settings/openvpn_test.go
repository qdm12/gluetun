package settings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ivpnAccountID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		s     string
		match bool
	}{
		{},
		{s: "abc"},
		{s: "i"},
		{s: "ivpn"},
		{s: "ivpn-aaaa"},
		{s: "ivpn-aaaa-aaaa"},
		{s: "ivpn-aaaa-aaaa-aaa"},
		{s: "ivpn-aaaa-aaaa-aaaa", match: true},
		{s: "ivpn-aaaa-aaaa-aaaaa"},
		{s: "ivpn-a6B7-fP91-Zh6Y", match: true},
		{s: "i-aaaa"},
		{s: "i-aaaa-aaaa"},
		{s: "i-aaaa-aaaa-aaa"},
		{s: "i-aaaa-aaaa-aaaa", match: true},
		{s: "i-aaaa-aaaa-aaaaa"},
		{s: "i-a6B7-fP91-Zh6Y", match: true},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.s, func(t *testing.T) {
			t.Parallel()

			match := ivpnAccountID.MatchString(testCase.s)

			assert.Equal(t, testCase.match, match)
		})
	}
}
