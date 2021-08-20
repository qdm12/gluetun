package wireguard

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vishvananda/netlink"
)

func Test_addRule(t *testing.T) {
	t.Parallel()

	rulePriority := 10000
	const firewallMark = 999

	cleanup, err := addRule(rulePriority, firewallMark)
	require.NoError(t, err)
	defer func() {
		err := cleanup()
		assert.NoError(t, err)
	}()

	rules, err := netlink.RuleList(netlink.FAMILY_ALL)
	require.NoError(t, err)
	var rule netlink.Rule
	var ruleFound bool
	for _, rule = range rules {
		if rule.Mark == firewallMark {
			ruleFound = true
			break
		}
	}
	require.True(t, ruleFound)
	expectedRule := netlink.Rule{
		Invert:            true,
		Priority:          rulePriority,
		Mark:              firewallMark,
		Table:             firewallMark,
		Mask:              4294967295,
		Goto:              -1,
		Flow:              -1,
		SuppressIfgroup:   -1,
		SuppressPrefixlen: -1,
	}
	assert.Equal(t, expectedRule, rule)

	// Existing rule cannot be added
	nilCleanup, err := addRule(rulePriority, firewallMark)
	if nilCleanup != nil {
		_ = nilCleanup() // in case it succeeds
	}
	require.Error(t, err)
	assert.Equal(t, "file exists", err.Error())
}
