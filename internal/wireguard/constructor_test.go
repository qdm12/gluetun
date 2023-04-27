package wireguard

import (
	"net"
	"net/netip"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	t.Parallel()

	const validKeyString = "oMNSf/zJ0pt1ciy+qIRk8Rlyfs9accwuRLnKd85Yl1Q="
	logger := NewMockLogger(nil)
	netLinker := NewMockNetLinker(nil)

	testCases := map[string]struct {
		settings  Settings
		wireguard *Wireguard
		err       error
	}{
		"bad settings": {
			settings: Settings{
				PrivateKey: "",
			},
			err: ErrPrivateKeyMissing,
		},
		"minimal valid settings": {
			settings: Settings{
				PrivateKey: validKeyString,
				PublicKey:  validKeyString,
				Endpoint: &net.UDPAddr{
					IP: net.IPv4(1, 2, 3, 4),
				},
				Addresses: []netip.Prefix{
					netip.PrefixFrom(netip.AddrFrom4([4]byte{5, 6, 7, 8}), 32),
				},
				FirewallMark: 100,
			},
			wireguard: &Wireguard{
				logger:  logger,
				netlink: netLinker,
				settings: Settings{
					InterfaceName: "wg0",
					PrivateKey:    validKeyString,
					PublicKey:     validKeyString,
					Endpoint: &net.UDPAddr{
						IP:   net.IPv4(1, 2, 3, 4),
						Port: 51820,
					},
					Addresses: []netip.Prefix{
						netip.PrefixFrom(netip.AddrFrom4([4]byte{5, 6, 7, 8}), 32),
					},
					FirewallMark:   100,
					IPv6:           ptr(false),
					Implementation: "auto",
				},
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			wireguard, err := New(testCase.settings, netLinker, logger)

			if testCase.err != nil {
				require.Error(t, err)
				assert.Equal(t, testCase.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.wireguard, wireguard)
		})
	}
}
