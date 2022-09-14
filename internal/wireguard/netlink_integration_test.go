//go:build netlink
// +build netlink

package wireguard

import (
	"net"
	"testing"

	"github.com/qdm12/gluetun/internal/netlink"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"
)

func Test_netlink_Wireguard_addAddresses(t *testing.T) {
	t.Parallel()

	netlinker := netlink.New()

	linkAttrs := netlink.NewLinkAttrs()
	linkAttrs.Name = "test_8081"
	link := &netlink.Bridge{
		LinkAttrs: linkAttrs,
	}

	// Remove any previously created test interface from a crashed/panic
	// test or test suite run.
	err := netlinker.LinkDel(link)
	if err != nil && err.Error() != "invalid argument" {
		require.NoError(t, err)
	}

	err = netlinker.LinkAdd(link)
	require.NoError(t, err)

	defer func() {
		err = netlinker.LinkDel(link)
		assert.NoError(t, err)
	}()

	addresses := []*net.IPNet{
		{IP: net.IP{1, 2, 3, 4}, Mask: net.IPv4Mask(255, 255, 255, 255)},
		{IP: net.IP{5, 6, 7, 8}, Mask: net.IPv4Mask(255, 255, 255, 255)},
	}

	wg := &Wireguard{
		netlink: netlinker,
	}

	// Success
	err = wg.addAddresses(link, addresses)
	require.NoError(t, err)

	netlinkAddresses, err := netlinker.AddrList(link, netlink.FAMILY_ALL)
	require.NoError(t, err)
	require.Equal(t, len(addresses), len(netlinkAddresses))
	for i, netlinkAddress := range netlinkAddresses {
		ipNet := netlinkAddress.IPNet
		assert.Equal(t, addresses[i], ipNet)
	}

	// Existing address cannot be added
	err = wg.addAddresses(link, addresses)
	require.Error(t, err)
	assert.EqualError(t, err, "file exists: when adding address 1.2.3.4/32 to link test_8081")
}

func Test_netlink_Wireguard_addRule(t *testing.T) {
	t.Parallel()

	netlinker := netlink.New()
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

	rules, err := netlinker.RuleList(netlink.FAMILY_V4)
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
	nilCleanup, err := wg.addRule(rulePriority,
		firewallMark, family)
	if nilCleanup != nil {
		_ = nilCleanup() // in case it succeeds
	}
	require.Error(t, err)
	assert.Equal(t, "cannot add rule ip rule 10000: from all to all table 999: file exists", err.Error())
}
