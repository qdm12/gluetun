package routing

import (
	"errors"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/qdm12/gluetun/internal/netlink"
	"github.com/qdm12/gluetun/internal/netlink/mock_netlink"
	"github.com/qdm12/golibs/logging/mock_logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeIPNet(t *testing.T, n byte) *net.IPNet {
	t.Helper()
	return &net.IPNet{
		IP:   net.IPv4(n, n, n, 0),
		Mask: net.IPv4Mask(255, 255, 255, 0),
	}
}

func makeIPRule(t *testing.T, src, dst *net.IPNet,
	table, priority int) *netlink.Rule {
	t.Helper()
	rule := netlink.NewRule()
	rule.Src = src
	rule.Dst = dst
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
		src      *net.IPNet
		dst      *net.IPNet
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
			err: errors.New("cannot list rules: dummy error"),
		},
		"rule already exists": {
			src:      makeIPNet(t, 1),
			dst:      makeIPNet(t, 2),
			table:    99,
			priority: 99,
			dbgMsg:   "ip rule add from 1.1.1.0/24 to 2.2.2.0/24 lookup 99 pref 99",
			ruleList: ruleListCall{
				rules: []netlink.Rule{
					*makeIPRule(t, makeIPNet(t, 2), makeIPNet(t, 2), 99, 99),
					*makeIPRule(t, makeIPNet(t, 1), makeIPNet(t, 2), 99, 99),
				},
			},
		},
		"add rule error": {
			src:      makeIPNet(t, 1),
			dst:      makeIPNet(t, 2),
			table:    99,
			priority: 99,
			dbgMsg:   "ip rule add from 1.1.1.0/24 to 2.2.2.0/24 lookup 99 pref 99",
			ruleAdd: ruleAddCall{
				expected:  true,
				ruleToAdd: makeIPRule(t, makeIPNet(t, 1), makeIPNet(t, 2), 99, 99),
				err:       errDummy,
			},
			err: errors.New("dummy error: for rule: ip rule 99: from 1.1.1.0/24 table 99"),
		},
		"add rule success": {
			src:      makeIPNet(t, 1),
			dst:      makeIPNet(t, 2),
			table:    99,
			priority: 99,
			dbgMsg:   "ip rule add from 1.1.1.0/24 to 2.2.2.0/24 lookup 99 pref 99",
			ruleList: ruleListCall{
				rules: []netlink.Rule{
					*makeIPRule(t, makeIPNet(t, 2), makeIPNet(t, 2), 99, 99),
					*makeIPRule(t, makeIPNet(t, 1), makeIPNet(t, 2), 101, 101),
				},
			},
			ruleAdd: ruleAddCall{
				expected:  true,
				ruleToAdd: makeIPRule(t, makeIPNet(t, 1), makeIPNet(t, 2), 99, 99),
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			logger := mock_logging.NewMockLogger(ctrl)
			logger.EXPECT().Debug(testCase.dbgMsg)

			netLinker := mock_netlink.NewMockNetLinker(ctrl)
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
		src      *net.IPNet
		dst      *net.IPNet
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
			err: errors.New("cannot list rules: dummy error"),
		},
		"rule delete error": {
			src:      makeIPNet(t, 1),
			dst:      makeIPNet(t, 2),
			table:    99,
			priority: 99,
			dbgMsg:   "ip rule del from 1.1.1.0/24 to 2.2.2.0/24 lookup 99 pref 99",
			ruleList: ruleListCall{
				rules: []netlink.Rule{
					*makeIPRule(t, makeIPNet(t, 1), makeIPNet(t, 2), 99, 99),
				},
			},
			ruleDel: ruleDelCall{
				expected:  true,
				ruleToDel: makeIPRule(t, makeIPNet(t, 1), makeIPNet(t, 2), 99, 99),
				err:       errDummy,
			},
			err: errors.New("dummy error: for rule: ip rule 99: from 1.1.1.0/24 table 99"),
		},
		"rule deleted": {
			src:      makeIPNet(t, 1),
			dst:      makeIPNet(t, 2),
			table:    99,
			priority: 99,
			dbgMsg:   "ip rule del from 1.1.1.0/24 to 2.2.2.0/24 lookup 99 pref 99",
			ruleList: ruleListCall{
				rules: []netlink.Rule{
					*makeIPRule(t, makeIPNet(t, 2), makeIPNet(t, 2), 99, 99),
					*makeIPRule(t, makeIPNet(t, 1), makeIPNet(t, 2), 99, 99),
				},
			},
			ruleDel: ruleDelCall{
				expected:  true,
				ruleToDel: makeIPRule(t, makeIPNet(t, 1), makeIPNet(t, 2), 99, 99),
			},
		},
		"rule does not exist": {
			src:      makeIPNet(t, 1),
			dst:      makeIPNet(t, 2),
			table:    99,
			priority: 99,
			dbgMsg:   "ip rule del from 1.1.1.0/24 to 2.2.2.0/24 lookup 99 pref 99",
			ruleList: ruleListCall{
				rules: []netlink.Rule{
					*makeIPRule(t, makeIPNet(t, 2), makeIPNet(t, 2), 99, 99),
					*makeIPRule(t, makeIPNet(t, 1), makeIPNet(t, 2), 101, 101),
				},
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			logger := mock_logging.NewMockLogger(ctrl)
			logger.EXPECT().Debug(testCase.dbgMsg)

			netLinker := mock_netlink.NewMockNetLinker(ctrl)
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
		src      *net.IPNet
		dst      *net.IPNet
		table    int
		priority int
		dbgMsg   string
	}{
		"default values": {
			dbgMsg: "ip rule del pref 0",
		},
		"add rule": {
			add:      true,
			src:      makeIPNet(t, 1),
			dst:      makeIPNet(t, 2),
			table:    100,
			priority: 101,
			dbgMsg:   "ip rule add from 1.1.1.0/24 to 2.2.2.0/24 lookup 100 pref 101",
		},
		"del rule": {
			src:      makeIPNet(t, 1),
			dst:      makeIPNet(t, 2),
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
