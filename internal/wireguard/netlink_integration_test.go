//go:build netlink && linux

package wireguard

import (
	"net/netip"
	"testing"

	"github.com/qdm12/gluetun/internal/netlink"
	"github.com/qdm12/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"
)

type noopDebugLogger struct{}

func (n noopDebugLogger) Debugf(format string, args ...any) {}
func (n noopDebugLogger) Patch(options ...log.Option)       {}

func Test_netlink_Wireguard_addAddresses(t *testing.T) {
	t.Parallel()

	netlinker := netlink.New(&noopDebugLogger{})

	link := netlink.Link{
		Type: "bridge",
		Name: "test_8081",
	}

	// Remove any previously created test interface from a crashed/panic
	// test or test suite run.
	err := netlinker.LinkDel(link)
	if err != nil && err.Error() != "invalid argument" {
		require.NoError(t, err)
	}

	linkIndex, err := netlinker.LinkAdd(link)
	require.NoError(t, err)
	link.Index = linkIndex

	defer func() {
		err = netlinker.LinkDel(link)
		assert.NoError(t, err)
	}()

	addresses := []netip.Prefix{
		netip.PrefixFrom(netip.AddrFrom4([4]byte{1, 2, 3, 4}), 32),
		netip.PrefixFrom(netip.AddrFrom4([4]byte{5, 6, 7, 8}), 32),
	}

	wg := &Wireguard{
		netlink: netlinker,
		settings: Settings{
			IPv6: new(bool),
		},
	}

	const addIterations = 2 // initial + replace

	for i := 0; i < addIterations; i++ {
		err = wg.addAddresses(link, addresses)
		require.NoError(t, err)

		netlinkAddresses, err := netlinker.AddrList(link, netlink.FamilyAll)
		require.NoError(t, err)
		require.Equal(t, len(addresses), len(netlinkAddresses))
		for i, netlinkAddress := range netlinkAddresses {
			require.NotNil(t, netlinkAddress.Network)
			assert.Equal(t, addresses[i], netlinkAddress.Network)
		}
	}
}

func Test_netlink_Wireguard_addRule(t *testing.T) {
	t.Parallel()

	netlinker := netlink.New(&noopDebugLogger{})
	wg := &Wireguard{
		netlink: netlinker,
	}

	rulePriority := 10000
	const firewallMark = 999
	const family = unix.AF_INET // ipv4

	cleanup, err := wg.addRule(rulePriority,
		firewallMark, family)
	require.NoError(t, err)
	defer func() {
		err := cleanup()
		assert.NoError(t, err)
	}()

	rules, err := netlinker.RuleList(netlink.FamilyV4)
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
		Invert:   true,
		Priority: rulePriority,
		Mark:     firewallMark,
		Table:    firewallMark,
	}
	assert.Equal(t, expectedRule, rule)

	// Existing rule cannot be added
	nilCleanup, err := wg.addRule(rulePriority,
		firewallMark, family)
	if nilCleanup != nil {
		_ = nilCleanup() // in case it succeeds
	}
	require.Error(t, err)
	assert.EqualError(t, err, "adding ip rule 10000: from all to all table 999: file exists")
}
