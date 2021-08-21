package wireguard

import (
	"errors"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vishvananda/netlink"
)

func Test_Wireguard_addAddresses(t *testing.T) {
	t.Parallel()

	ipNetOne := &net.IPNet{IP: net.IPv4(1, 2, 3, 4), Mask: net.IPv4Mask(255, 255, 255, 255)}
	ipNetTwo := &net.IPNet{IP: net.IPv4(4, 5, 6, 7), Mask: net.IPv4Mask(255, 255, 255, 128)}

	newLink := func() netlink.Link {
		linkAttrs := netlink.NewLinkAttrs()
		linkAttrs.Name = "a_bridge"
		return &netlink.Bridge{
			LinkAttrs: linkAttrs,
		}
	}

	errDummy := errors.New("dummy")

	testCases := map[string]struct {
		link          netlink.Link
		addrs         []*net.IPNet
		expectedAddrs []*netlink.Addr
		addrAddErrs   []error
		err           error
	}{
		"success": {
			link:  newLink(),
			addrs: []*net.IPNet{ipNetOne, ipNetTwo},
			expectedAddrs: []*netlink.Addr{
				{IPNet: ipNetOne}, {IPNet: ipNetTwo},
			},
			addrAddErrs: []error{nil, nil},
		},
		"first add error": {
			link:  newLink(),
			addrs: []*net.IPNet{ipNetOne, ipNetTwo},
			expectedAddrs: []*netlink.Addr{
				{IPNet: ipNetOne},
			},
			addrAddErrs: []error{errDummy},
			err:         errors.New("dummy: when adding address 1.2.3.4/32 to link a_bridge"),
		},
		"second add error": {
			link:  newLink(),
			addrs: []*net.IPNet{ipNetOne, ipNetTwo},
			expectedAddrs: []*netlink.Addr{
				{IPNet: ipNetOne}, {IPNet: ipNetTwo},
			},
			addrAddErrs: []error{nil, errDummy},
			err:         errors.New("dummy: when adding address 4.5.6.7/25 to link a_bridge"),
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			require.Equal(t, len(testCase.expectedAddrs), len(testCase.addrAddErrs))

			netLinker := NewMockNetLinker(ctrl)
			wg := Wireguard{
				netlink: netLinker,
			}

			for i := range testCase.expectedAddrs {
				netLinker.EXPECT().
					AddrAdd(testCase.link, testCase.expectedAddrs[i]).
					Return(testCase.addrAddErrs[i])
			}

			err := wg.addAddresses(testCase.link, testCase.addrs)

			if testCase.err != nil {
				require.Error(t, err)
				assert.Equal(t, testCase.err.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
