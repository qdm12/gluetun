package wireguard

import (
	"net"
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
				Addresses: []*net.IPNet{{
					IP:   net.IPv4(5, 6, 7, 8),
					Mask: net.IPv4Mask(255, 255, 255, 255)},
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
					Addresses: []*net.IPNet{{
						IP:   net.IPv4(5, 6, 7, 8),
						Mask: net.IPv4Mask(255, 255, 255, 255)},
					},
					FirewallMark: 100,
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
