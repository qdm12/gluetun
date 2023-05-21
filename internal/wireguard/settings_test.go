package wireguard

import (
	"errors"
	"net/netip"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.zx2c4.com/wireguard/device"
)

func ptr[T any](v T) *T { return &v }

func Test_Settings_SetDefaults(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		original Settings
		expected Settings
	}{
		"empty settings": {
			expected: Settings{
				InterfaceName:  "wg0",
				FirewallMark:   51820,
				MTU:            device.DefaultMTU,
				IPv6:           ptr(false),
				Implementation: "auto",
			},
		},
		"default endpoint port": {
			original: Settings{
				Endpoint: netip.AddrPortFrom(netip.AddrFrom4([4]byte{1, 2, 3, 4}), 0),
			},
			expected: Settings{
				InterfaceName:  "wg0",
				FirewallMark:   51820,
				Endpoint:       netip.AddrPortFrom(netip.AddrFrom4([4]byte{1, 2, 3, 4}), 51820),
				MTU:            device.DefaultMTU,
				IPv6:           ptr(false),
				Implementation: "auto",
			},
		},
		"not empty settings": {
			original: Settings{
				InterfaceName:  "wg1",
				FirewallMark:   999,
				Endpoint:       netip.AddrPortFrom(netip.AddrFrom4([4]byte{1, 2, 3, 4}), 9999),
				MTU:            device.DefaultMTU,
				IPv6:           ptr(true),
				Implementation: "userspace",
			},
			expected: Settings{
				InterfaceName:  "wg1",
				FirewallMark:   999,
				Endpoint:       netip.AddrPortFrom(netip.AddrFrom4([4]byte{1, 2, 3, 4}), 9999),
				MTU:            device.DefaultMTU,
				IPv6:           ptr(true),
				Implementation: "userspace",
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			testCase.original.SetDefaults()

			assert.Equal(t, testCase.expected, testCase.original)
		})
	}
}

func Test_Settings_Check(t *testing.T) {
	t.Parallel()

	const (
		validKey1 = "oMNSf/zJ0pt1ciy+qIRk8Rlyfs9accwuRLnKd85Yl1Q="
		validKey2 = "aPjc9US5ICB30D1P4glR9tO7bkB2Ga+KZiFqnoypBHk="
	)

	testCases := map[string]struct {
		settings Settings
		err      error
	}{
		"empty settings": {
			err: errors.New("invalid interface name: "),
		},
		"bad interface name": {
			settings: Settings{
				InterfaceName: "$H1T",
			},
			err: errors.New("invalid interface name: $H1T"),
		},
		"empty private key": {
			settings: Settings{
				InterfaceName: "wg0",
			},
			err: ErrPrivateKeyMissing,
		},
		"bad private key": {
			settings: Settings{
				InterfaceName: "wg0",
				PrivateKey:    "bad key",
			},
			err: ErrPrivateKeyInvalid,
		},
		"empty public key": {
			settings: Settings{
				InterfaceName: "wg0",
				PrivateKey:    validKey1,
			},
			err: ErrPublicKeyMissing,
		},
		"bad public key": {
			settings: Settings{
				InterfaceName: "wg0",
				PrivateKey:    validKey1,
				PublicKey:     "bad key",
			},
			err: errors.New("cannot parse public key: bad key"),
		},
		"bad preshared key": {
			settings: Settings{
				InterfaceName: "wg0",
				PrivateKey:    validKey1,
				PublicKey:     validKey2,
				PreSharedKey:  "bad key",
			},
			err: errors.New("cannot parse pre-shared key"),
		},
		"invalid endpoint address": {
			settings: Settings{
				InterfaceName: "wg0",
				PrivateKey:    validKey1,
				PublicKey:     validKey2,
			},
			err: ErrEndpointAddrMissing,
		},
		"zero endpoint port": {
			settings: Settings{
				InterfaceName: "wg0",
				PrivateKey:    validKey1,
				PublicKey:     validKey2,
				Endpoint:      netip.AddrPortFrom(netip.AddrFrom4([4]byte{1, 2, 3, 4}), 0),
			},
			err: ErrEndpointPortMissing,
		},
		"no address": {
			settings: Settings{
				InterfaceName: "wg0",
				PrivateKey:    validKey1,
				PublicKey:     validKey2,
				Endpoint:      netip.AddrPortFrom(netip.AddrFrom4([4]byte{1, 2, 3, 4}), 51820),
			},
			err: ErrAddressMissing,
		},
		"invalid address": {
			settings: Settings{
				InterfaceName: "wg0",
				PrivateKey:    validKey1,
				PublicKey:     validKey2,
				Endpoint:      netip.AddrPortFrom(netip.AddrFrom4([4]byte{1, 2, 3, 4}), 51820),
				Addresses:     []netip.Prefix{{}},
			},
			err: errors.New("interface address is not valid: for address 1 of 1"),
		},
		"zero firewall mark": {
			settings: Settings{
				InterfaceName: "wg0",
				PrivateKey:    validKey1,
				PublicKey:     validKey2,
				Endpoint:      netip.AddrPortFrom(netip.AddrFrom4([4]byte{1, 2, 3, 4}), 51820),
				Addresses: []netip.Prefix{
					netip.PrefixFrom(netip.AddrFrom4([4]byte{1, 2, 3, 4}), 24),
				},
			},
			err: ErrFirewallMarkMissing,
		},
		"missing_MTU": {
			settings: Settings{
				InterfaceName: "wg0",
				PrivateKey:    validKey1,
				PublicKey:     validKey2,
				Endpoint:      netip.AddrPortFrom(netip.AddrFrom4([4]byte{1, 2, 3, 4}), 51820),
				Addresses: []netip.Prefix{
					netip.PrefixFrom(netip.AddrFrom4([4]byte{1, 2, 3, 4}), 24),
				},
				FirewallMark: 999,
			},
			err: ErrMTUMissing,
		},
		"invalid implementation": {
			settings: Settings{
				InterfaceName: "wg0",
				PrivateKey:    validKey1,
				PublicKey:     validKey2,
				Endpoint:      netip.AddrPortFrom(netip.AddrFrom4([4]byte{1, 2, 3, 4}), 51820),
				Addresses: []netip.Prefix{
					netip.PrefixFrom(netip.AddrFrom4([4]byte{1, 2, 3, 4}), 24),
				},
				FirewallMark:   999,
				MTU:            1420,
				Implementation: "x",
			},
			err: errors.New("invalid implementation: x"),
		},
		"all valid": {
			settings: Settings{
				InterfaceName: "wg0",
				PrivateKey:    validKey1,
				PublicKey:     validKey2,
				Endpoint:      netip.AddrPortFrom(netip.AddrFrom4([4]byte{1, 2, 3, 4}), 51820),
				Addresses: []netip.Prefix{
					netip.PrefixFrom(netip.AddrFrom4([4]byte{1, 2, 3, 4}), 24),
				},
				FirewallMark:   999,
				MTU:            1420,
				Implementation: "userspace",
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := testCase.settings.Check()

			if testCase.err != nil {
				require.Error(t, err)
				assert.Equal(t, testCase.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func toStringPtr(s string) *string { return &s }

func Test_ToLinesSettings_setDefaults(t *testing.T) {
	t.Parallel()

	settings := ToLinesSettings{
		Indent: toStringPtr("indent"),
	}

	someFunc := func(settings ToLinesSettings) {
		settings.setDefaults()
		expectedSettings := ToLinesSettings{
			Indent:          toStringPtr("indent"),
			FieldPrefix:     toStringPtr("├── "),
			LastFieldPrefix: toStringPtr("└── "),
		}
		assert.Equal(t, expectedSettings, settings)
	}
	someFunc(settings)

	untouchedSettings := ToLinesSettings{
		Indent: toStringPtr("indent"),
	}
	assert.Equal(t, untouchedSettings, settings)
}

func Test_Settings_String(t *testing.T) {
	t.Parallel()

	settings := Settings{
		InterfaceName:  "wg0",
		IPv6:           ptr(true),
		Implementation: "x",
	}
	const expected = `├── Interface name: wg0
├── Private key: not set
├── Pre shared key: not set
├── Endpoint: not set
├── IPv6: enabled
├── Implementation: x
└── Addresses: not set`
	s := settings.String()
	assert.Equal(t, expected, s)
}

func Test_Settings_Lines(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		settings     Settings
		lineSettings ToLinesSettings
		lines        []string
	}{
		"empty settings": {
			settings: Settings{
				IPv6: ptr(false),
			},
			lines: []string{
				"├── Interface name: ",
				"├── Private key: not set",
				"├── Pre shared key: not set",
				"├── Endpoint: not set",
				"├── IPv6: disabled",
				"├── Implementation: ",
				"└── Addresses: not set",
			},
		},
		"settings all set": {
			settings: Settings{
				InterfaceName: "wg0",
				PrivateKey:    "private key",
				PublicKey:     "public key",
				PreSharedKey:  "pre-shared key",
				Endpoint:      netip.AddrPortFrom(netip.AddrFrom4([4]byte{1, 2, 3, 4}), 51820),
				FirewallMark:  999,
				RulePriority:  888,
				Addresses: []netip.Prefix{
					netip.PrefixFrom(netip.AddrFrom4([4]byte{1, 1, 1, 1}), 24),
					netip.PrefixFrom(netip.AddrFrom4([4]byte{2, 2, 2, 2}), 32),
				},
				IPv6:           ptr(true),
				Implementation: "userspace",
			},
			lines: []string{
				"├── Interface name: wg0",
				"├── Private key: set",
				"├── PublicKey: public key",
				"├── Pre shared key: set",
				"├── Endpoint: 1.2.3.4:51820",
				"├── IPv6: enabled",
				"├── Firewall mark: 999",
				"├── Rule priority: 888",
				"├── Implementation: userspace",
				"└── Addresses:",
				"    ├── 1.1.1.1/24",
				"    └── 2.2.2.2/32",
			},
		},
		"custom line settings": {
			lineSettings: ToLinesSettings{
				Indent:          toStringPtr("  "),
				FieldPrefix:     toStringPtr("- "),
				LastFieldPrefix: toStringPtr("* "),
			},
			settings: Settings{
				InterfaceName: "wg0",
				Addresses: []netip.Prefix{
					netip.PrefixFrom(netip.AddrFrom4([4]byte{1, 1, 1, 1}), 24),
					netip.PrefixFrom(netip.AddrFrom4([4]byte{2, 2, 2, 2}), 32),
				},
				IPv6: ptr(false),
			},
			lines: []string{
				"- Interface name: wg0",
				"- Private key: not set",
				"- Pre shared key: not set",
				"- Endpoint: not set",
				"- IPv6: disabled",
				"- Implementation: ",
				"* Addresses:",
				"  - 1.1.1.1/24",
				"  * 2.2.2.2/32",
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			lines := testCase.settings.ToLines(testCase.lineSettings)

			assert.Equal(t, testCase.lines, lines)
		})
	}
}
