package wireguard

import (
	"errors"
	"net/netip"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Wireguard_addAddresses(t *testing.T) {
	t.Parallel()

	ipNetOne := netip.PrefixFrom(netip.AddrFrom4([4]byte{1, 2, 3, 4}), 32)
	ipNetTwo := netip.PrefixFrom(netip.MustParseAddr("::1234"), 64)

	errDummy := errors.New("dummy")

	testCases := map[string]struct {
		linkIndex uint32
		addrs     []netip.Prefix
		wgBuilder func(ctrl *gomock.Controller, linkIndex uint32) *Wireguard
		err       error
	}{
		"success": {
			linkIndex: 1,
			addrs:     []netip.Prefix{ipNetOne, ipNetTwo},
			wgBuilder: func(ctrl *gomock.Controller, linkIndex uint32) *Wireguard {
				netLinker := NewMockNetLinker(ctrl)
				firstCall := netLinker.EXPECT().
					AddrReplace(linkIndex, ipNetOne).
					Return(nil)
				netLinker.EXPECT().
					AddrReplace(linkIndex, ipNetTwo).
					Return(nil).After(firstCall)
				return &Wireguard{
					netlink: netLinker,
					settings: Settings{
						IPv6: ptrTo(true),
					},
				}
			},
		},
		"first add error": {
			linkIndex: 1,
			addrs:     []netip.Prefix{ipNetOne, ipNetTwo},
			wgBuilder: func(ctrl *gomock.Controller, linkIndex uint32) *Wireguard {
				netLinker := NewMockNetLinker(ctrl)
				netLinker.EXPECT().
					AddrReplace(linkIndex, ipNetOne).
					Return(errDummy)
				return &Wireguard{
					netlink: netLinker,
					settings: Settings{
						IPv6: ptrTo(true),
					},
				}
			},
			err: errors.New("dummy: when adding address 1.2.3.4/32 to link with index 1"),
		},
		"second add error": {
			linkIndex: 1,
			addrs:     []netip.Prefix{ipNetOne, ipNetTwo},
			wgBuilder: func(ctrl *gomock.Controller, linkIndex uint32) *Wireguard {
				netLinker := NewMockNetLinker(ctrl)
				firstCall := netLinker.EXPECT().
					AddrReplace(linkIndex, ipNetOne).
					Return(nil)
				netLinker.EXPECT().
					AddrReplace(linkIndex, ipNetTwo).
					Return(errDummy).After(firstCall)
				return &Wireguard{
					netlink: netLinker,
					settings: Settings{
						IPv6: ptrTo(true),
					},
				}
			},
			err: errors.New("dummy: when adding address ::1234/64 to link with index 1"),
		},
		"ignore IPv6": {
			addrs: []netip.Prefix{ipNetTwo},
			wgBuilder: func(_ *gomock.Controller, _ uint32) *Wireguard {
				return &Wireguard{
					settings: Settings{
						IPv6: ptrTo(false),
					},
				}
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			wg := testCase.wgBuilder(ctrl, testCase.linkIndex)

			err := wg.addAddresses(testCase.linkIndex, testCase.addrs)

			if testCase.err != nil {
				require.Error(t, err)
				assert.Equal(t, testCase.err.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
