package netlink

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ruleDbgMsg(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		add    bool
		rule   Rule
		dbgMsg string
	}{
		"default values": {
			dbgMsg: "ip -f 0 rule del pref 0",
		},
		"add rule": {
			add: true,
			rule: Rule{
				Family:   FamilyV4,
				Src:      makeNetipPrefix(1),
				Dst:      makeNetipPrefix(2),
				Table:    100,
				Priority: 101,
			},
			dbgMsg: "ip -f inet rule add from 1.1.1.0/24 to 2.2.2.0/24 lookup 100 pref 101",
		},
		"del rule": {
			rule: Rule{
				Family:   FamilyV4,
				Src:      makeNetipPrefix(1),
				Dst:      makeNetipPrefix(2),
				Table:    100,
				Priority: 101,
			},
			dbgMsg: "ip -f inet rule del from 1.1.1.0/24 to 2.2.2.0/24 lookup 100 pref 101",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			dbgMsg := ruleDbgMsg(testCase.add, testCase.rule)

			assert.Equal(t, testCase.dbgMsg, dbgMsg)
		})
	}
}
