package wireguard

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/qdm12/gluetun/internal/netlink"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"
)

func Test_Wireguard_addRule(t *testing.T) {
	t.Parallel()

	const rulePriority = 987
	const firewallMark = 456
	const family = unix.AF_INET

	errDummy := errors.New("dummy")

	testCases := map[string]struct {
		expectedRule *netlink.Rule
		ruleAddErr   error
		err          error
		ruleDelErr   error
		cleanupErr   error
	}{
		"success": {
			expectedRule: &netlink.Rule{
				Invert:            true,
				Priority:          rulePriority,
				Mark:              firewallMark,
				Table:             firewallMark,
				Mask:              -1,
				Goto:              -1,
				Flow:              -1,
				SuppressIfgroup:   -1,
				SuppressPrefixlen: -1,
				Family:            family,
			},
		},
		"rule add error": {
			expectedRule: &netlink.Rule{
				Invert:            true,
				Priority:          rulePriority,
				Mark:              firewallMark,
				Table:             firewallMark,
				Mask:              -1,
				Goto:              -1,
				Flow:              -1,
				SuppressIfgroup:   -1,
				SuppressPrefixlen: -1,
				Family:            family,
			},
			ruleAddErr: errDummy,
			err:        errors.New("cannot add rule ip rule 987: from all to all table 456: dummy"),
		},
		"rule delete error": {
			expectedRule: &netlink.Rule{
				Invert:            true,
				Priority:          rulePriority,
				Mark:              firewallMark,
				Table:             firewallMark,
				Mask:              -1,
				Goto:              -1,
				Flow:              -1,
				SuppressIfgroup:   -1,
				SuppressPrefixlen: -1,
				Family:            family,
			},
			ruleDelErr: errDummy,
			cleanupErr: errors.New("cannot delete rule ip rule 987: from all to all table 456: dummy"),
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			netLinker := NewMockNetLinker(ctrl)
			wg := Wireguard{
				netlink: netLinker,
			}

			netLinker.EXPECT().RuleAdd(testCase.expectedRule).
				Return(testCase.ruleAddErr)
			cleanup, err := wg.addRule(rulePriority, firewallMark, family)
			if testCase.err != nil {
				require.Error(t, err)
				assert.Equal(t, testCase.err.Error(), err.Error())
				return
			}

			require.NoError(t, err)

			netLinker.EXPECT().RuleDel(testCase.expectedRule).
				Return(testCase.ruleDelErr)
			err = cleanup()
			if testCase.cleanupErr != nil {
				require.Error(t, err)
				assert.Equal(t, testCase.cleanupErr.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
