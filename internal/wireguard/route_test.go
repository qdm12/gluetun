package wireguard

import (
	"errors"
	"net/netip"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/qdm12/gluetun/internal/netlink"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Wireguard_addRoute(t *testing.T) {
	t.Parallel()

	const linkIndex = 88

	ipPrefix := netip.PrefixFrom(netip.AddrFrom4([4]byte{1, 2, 3, 4}), 32)

	const firewallMark = 51820

	errDummy := errors.New("dummy")

	testCases := map[string]struct {
		link          netlink.Link
		dst           netip.Prefix
		expectedRoute netlink.Route
		routeAddErr   error
		err           error
	}{
		"success": {
			link: netlink.Link{
				Index: linkIndex,
			},
			dst: ipPrefix,
			expectedRoute: netlink.Route{
				LinkIndex: linkIndex,
				Dst:       ipPrefix,
				Table:     firewallMark,
			},
		},
		"route add error": {
			link: netlink.Link{
				Name:  "a_bridge",
				Index: linkIndex,
			},
			dst: ipPrefix,
			expectedRoute: netlink.Route{
				LinkIndex: linkIndex,
				Dst:       ipPrefix,
				Table:     firewallMark,
			},
			routeAddErr: errDummy,
			err:         errors.New("adding route for link a_bridge, destination 1.2.3.4/32 and table 51820: dummy"), //nolint:lll
		},
	}

	for name, testCase := range testCases {
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
