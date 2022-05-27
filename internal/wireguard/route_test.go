package wireguard

import (
	"errors"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/qdm12/gluetun/internal/netlink"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Wireguard_addRoute(t *testing.T) {
	t.Parallel()

	const linkIndex = 88
	newLink := func() netlink.Link {
		linkAttrs := netlink.NewLinkAttrs()
		linkAttrs.Name = "a_bridge"
		linkAttrs.Index = linkIndex
		return &netlink.Bridge{
			LinkAttrs: linkAttrs,
		}
	}
	ipNet := &net.IPNet{IP: net.IPv4(1, 2, 3, 4), Mask: net.IPv4Mask(255, 255, 255, 255)}
	const firewallMark = 51820

	errDummy := errors.New("dummy")

	testCases := map[string]struct {
		link          netlink.Link
		dst           *net.IPNet
		expectedRoute *netlink.Route
		routeAddErr   error
		err           error
	}{
		"success": {
			link: newLink(),
			dst:  ipNet,
			expectedRoute: &netlink.Route{
				LinkIndex: linkIndex,
				Dst:       ipNet,
				Table:     firewallMark,
			},
		},
		"route add error": {
			link: newLink(),
			dst:  ipNet,
			expectedRoute: &netlink.Route{
				LinkIndex: linkIndex,
				Dst:       ipNet,
				Table:     firewallMark,
			},
			routeAddErr: errDummy,
			err:         errors.New("cannot add route for link a_bridge, destination 1.2.3.4/32 and table 51820: dummy"), //nolint:lll
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

			netLinker.EXPECT().
				RouteAdd(testCase.expectedRoute).
				Return(testCase.routeAddErr)

			err := wg.addRoute(testCase.link, testCase.dst, firewallMark)

			if testCase.err != nil {
				require.Error(t, err)
				assert.Equal(t, testCase.err.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
