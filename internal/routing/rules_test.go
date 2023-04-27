package routing

import (
	"errors"
	"net"
	"net/netip"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/qdm12/gluetun/internal/netlink"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeNetipPrefix(n byte) *netip.Prefix {
	const bits = 24
	prefix := netip.PrefixFrom(netip.AddrFrom4([4]byte{n, n, n, 0}), bits)
	return &prefix
}

func makeIPRule(src, dst *netip.Prefix,
	table, priority int) *netlink.Rule {
	rule := netlink.NewRule()
	rule.Src = NetipPrefixToIPNet(src)
	rule.Dst = NetipPrefixToIPNet(dst)
	rule.Table = table
	rule.Priority = priority
	return rule
}

func Test_Routing_addIPRule(t *testing.T) {
	t.Parallel()

	errDummy := errors.New("dummy error")

	type ruleListCall struct {
		rules []netlink.Rule
		err   error
	}

	type ruleAddCall struct {
		expected  bool
		ruleToAdd *netlink.Rule
		err       error
	}

	testCases := map[string]struct {
		src      *netip.Prefix
		dst      *netip.Prefix
		table    int
		priority int
		dbgMsg   string
		ruleList ruleListCall
		ruleAdd  ruleAddCall
		err      error
	}{
		"list error": {
			dbgMsg: "ip rule add pref 0",
			ruleList: ruleListCall{
				err: errDummy,
			},
			err: errors.New("listing rules: dummy error"),
		},
		"rule already exists": {
			src:      makeNetipPrefix(1),
			dst:      makeNetipPrefix(2),
			table:    99,
			priority: 99,
			dbgMsg:   "ip rule add from 1.1.1.0/24 to 2.2.2.0/24 lookup 99 pref 99",
			ruleList: ruleListCall{
				rules: []netlink.Rule{
					*makeIPRule(makeNetipPrefix(2), makeNetipPrefix(2), 99, 99),
					*makeIPRule(makeNetipPrefix(1), makeNetipPrefix(2), 99, 99),
				},
			},
		},
		"add rule error": {
			src:      makeNetipPrefix(1),
			dst:      makeNetipPrefix(2),
			table:    99,
			priority: 99,
			dbgMsg:   "ip rule add from 1.1.1.0/24 to 2.2.2.0/24 lookup 99 pref 99",
			ruleAdd: ruleAddCall{
				expected:  true,
				ruleToAdd: makeIPRule(makeNetipPrefix(1), makeNetipPrefix(2), 99, 99),
				err:       errDummy,
			},
			err: errors.New("adding rule ip rule 99: from 1.1.1.0/24 to 2.2.2.0/24 table 99: dummy error"),
		},
		"add rule success": {
			src:      makeNetipPrefix(1),
			dst:      makeNetipPrefix(2),
			table:    99,
			priority: 99,
			dbgMsg:   "ip rule add from 1.1.1.0/24 to 2.2.2.0/24 lookup 99 pref 99",
			ruleList: ruleListCall{
				rules: []netlink.Rule{
					*makeIPRule(makeNetipPrefix(2), makeNetipPrefix(2), 99, 99),
					*makeIPRule(makeNetipPrefix(1), makeNetipPrefix(2), 101, 101),
				},
			},
			ruleAdd: ruleAddCall{
				expected:  true,
				ruleToAdd: makeIPRule(makeNetipPrefix(1), makeNetipPrefix(2), 99, 99),
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			logger := NewMockLogger(ctrl)
			logger.EXPECT().Debug(testCase.dbgMsg)

			netLinker := NewMockNetLinker(ctrl)
			netLinker.EXPECT().RuleList(netlink.FAMILY_ALL).
				Return(testCase.ruleList.rules, testCase.ruleList.err)
			if testCase.ruleAdd.expected {
				netLinker.EXPECT().RuleAdd(testCase.ruleAdd.ruleToAdd).
					Return(testCase.ruleAdd.err)
			}

			r := Routing{
				logger:    logger,
				netLinker: netLinker,
			}

			err := r.addIPRule(testCase.src, testCase.dst,
				testCase.table, testCase.priority)

			if testCase.err != nil {
				require.Error(t, err)
				assert.Equal(t, testCase.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_Routing_deleteIPRule(t *testing.T) {
	t.Parallel()

	errDummy := errors.New("dummy error")

	type ruleListCall struct {
		rules []netlink.Rule
		err   error
	}

	type ruleDelCall struct {
		expected  bool
		ruleToDel *netlink.Rule
		err       error
	}

	testCases := map[string]struct {
		src      *netip.Prefix
		dst      *netip.Prefix
		table    int
		priority int
		dbgMsg   string
		ruleList ruleListCall
		ruleDel  ruleDelCall
		err      error
	}{
		"list error": {
			dbgMsg: "ip rule del pref 0",
			ruleList: ruleListCall{
				err: errDummy,
			},
			err: errors.New("listing rules: dummy error"),
		},
		"rule delete error": {
			src:      makeNetipPrefix(1),
			dst:      makeNetipPrefix(2),
			table:    99,
			priority: 99,
			dbgMsg:   "ip rule del from 1.1.1.0/24 to 2.2.2.0/24 lookup 99 pref 99",
			ruleList: ruleListCall{
				rules: []netlink.Rule{
					*makeIPRule(makeNetipPrefix(1), makeNetipPrefix(2), 99, 99),
				},
			},
			ruleDel: ruleDelCall{
				expected:  true,
				ruleToDel: makeIPRule(makeNetipPrefix(1), makeNetipPrefix(2), 99, 99),
				err:       errDummy,
			},
			err: errors.New("deleting rule ip rule 99: from 1.1.1.0/24 to 2.2.2.0/24 table 99: dummy error"),
		},
		"rule deleted": {
			src:      makeNetipPrefix(1),
			dst:      makeNetipPrefix(2),
			table:    99,
			priority: 99,
			dbgMsg:   "ip rule del from 1.1.1.0/24 to 2.2.2.0/24 lookup 99 pref 99",
			ruleList: ruleListCall{
				rules: []netlink.Rule{
					*makeIPRule(makeNetipPrefix(2), makeNetipPrefix(2), 99, 99),
					*makeIPRule(makeNetipPrefix(1), makeNetipPrefix(2), 99, 99),
				},
			},
			ruleDel: ruleDelCall{
				expected:  true,
				ruleToDel: makeIPRule(makeNetipPrefix(1), makeNetipPrefix(2), 99, 99),
			},
		},
		"rule does not exist": {
			src:      makeNetipPrefix(1),
			dst:      makeNetipPrefix(2),
			table:    99,
			priority: 99,
			dbgMsg:   "ip rule del from 1.1.1.0/24 to 2.2.2.0/24 lookup 99 pref 99",
			ruleList: ruleListCall{
				rules: []netlink.Rule{
					*makeIPRule(makeNetipPrefix(2), makeNetipPrefix(2), 99, 99),
					*makeIPRule(makeNetipPrefix(1), makeNetipPrefix(2), 101, 101),
				},
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			logger := NewMockLogger(ctrl)
			logger.EXPECT().Debug(testCase.dbgMsg)

			netLinker := NewMockNetLinker(ctrl)
			netLinker.EXPECT().RuleList(netlink.FAMILY_ALL).
				Return(testCase.ruleList.rules, testCase.ruleList.err)
			if testCase.ruleDel.expected {
				netLinker.EXPECT().RuleDel(testCase.ruleDel.ruleToDel).
					Return(testCase.ruleDel.err)
			}

			r := Routing{
				logger:    logger,
				netLinker: netLinker,
			}

			err := r.deleteIPRule(testCase.src, testCase.dst,
				testCase.table, testCase.priority)

			if testCase.err != nil {
				require.Error(t, err)
				assert.Equal(t, testCase.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_ruleDbgMsg(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		add      bool
		src      *netip.Prefix
		dst      *netip.Prefix
		table    int
		priority int
		dbgMsg   string
	}{
		"default values": {
			dbgMsg: "ip rule del pref 0",
		},
		"add rule": {
			add:      true,
			src:      makeNetipPrefix(1),
			dst:      makeNetipPrefix(2),
			table:    100,
			priority: 101,
			dbgMsg:   "ip rule add from 1.1.1.0/24 to 2.2.2.0/24 lookup 100 pref 101",
		},
		"del rule": {
			src:      makeNetipPrefix(1),
			dst:      makeNetipPrefix(2),
			table:    100,
			priority: 101,
			dbgMsg:   "ip rule del from 1.1.1.0/24 to 2.2.2.0/24 lookup 100 pref 101",
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			dbgMsg := ruleDbgMsg(testCase.add, testCase.src,
				testCase.dst, testCase.table, testCase.priority)

			assert.Equal(t, testCase.dbgMsg, dbgMsg)
		})
	}
}

func Test_rulesAreEqual(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		a     *netlink.Rule
		b     *netlink.Rule
		equal bool
	}{
		"both nil": {
			equal: true,
		},
		"first nil": {
			b: &netlink.Rule{},
		},
		"second nil": {
			a: &netlink.Rule{},
		},
		"both not nil": {
			a:     &netlink.Rule{},
			b:     &netlink.Rule{},
			equal: true,
		},
		"both equal": {
			a: &netlink.Rule{
				Src: &net.IPNet{
					IP:   net.IPv4(1, 1, 1, 1),
					Mask: net.IPv4Mask(255, 255, 255, 0),
				},
				Priority: 100,
				Table:    101,
			},
			b: &netlink.Rule{
				Src: &net.IPNet{
					IP:   net.IPv4(1, 1, 1, 1),
					Mask: net.IPv4Mask(255, 255, 255, 0),
				},
				Priority: 100,
				Table:    101,
			},
			equal: true,
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			equal := rulesAreEqual(testCase.a, testCase.b)

			assert.Equal(t, testCase.equal, equal)
		})
	}
}

func Test_ipNetsAreEqual(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		a     *net.IPNet
		b     *net.IPNet
		equal bool
	}{
		"both nil": {
			equal: true,
		},
		"first nil": {
			b: &net.IPNet{},
		},
		"second nil": {
			a: &net.IPNet{},
		},
		"both not nil": {
			a:     &net.IPNet{},
			b:     &net.IPNet{},
			equal: true,
		},
		"both equal": {
			a: &net.IPNet{
				IP:   net.IPv4(1, 1, 1, 1),
				Mask: net.IPv4Mask(255, 255, 255, 0),
			},
			b: &net.IPNet{
				IP:   net.IPv4(1, 1, 1, 1),
				Mask: net.IPv4Mask(255, 255, 255, 0),
			},
			equal: true,
		},
		"both not equal by IP": {
			a: &net.IPNet{
				IP:   net.IPv4(1, 1, 1, 1),
				Mask: net.IPv4Mask(255, 255, 255, 0),
			},
			b: &net.IPNet{
				IP:   net.IPv4(2, 2, 2, 2),
				Mask: net.IPv4Mask(255, 255, 255, 0),
			},
		},
		"both not equal by mask": {
			a: &net.IPNet{
				IP:   net.IPv4(1, 1, 1, 1),
				Mask: net.IPv4Mask(255, 255, 255, 255),
			},
			b: &net.IPNet{
				IP:   net.IPv4(1, 1, 1, 1),
				Mask: net.IPv4Mask(255, 255, 0, 0),
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			equal := ipNetsAreEqual(testCase.a, testCase.b)

			assert.Equal(t, testCase.equal, equal)
		})
	}
}
