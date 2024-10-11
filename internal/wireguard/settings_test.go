package wireguard

import (
	"net/netip"
	"testing"

	"github.com/stretchr/testify/assert"
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
				AllowedIPs:     []netip.Prefix{allIPv4()},
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
				AllowedIPs:     []netip.Prefix{allIPv4()},
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
				AllowedIPs:     []netip.Prefix{allIPv4()},
				MTU:            device.DefaultMTU,
				IPv6:           ptr(true),
				Implementation: "userspace",
			},
			expected: Settings{
				InterfaceName:  "wg1",
				FirewallMark:   999,
				Endpoint:       netip.AddrPortFrom(netip.AddrFrom4([4]byte{1, 2, 3, 4}), 9999),
				AllowedIPs:     []netip.Prefix{allIPv4()},
				MTU:            device.DefaultMTU,
				IPv6:           ptr(true),
				Implementation: "userspace",
			},
		},
	}

	for name, testCase := range testCases {
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
		settings   Settings
		errWrapped error
		errMessage string
	}{
		"empty settings": {
			errWrapped: ErrInterfaceNameInvalid,
			errMessage: "invalid interface name: ",
		},
		"bad interface name": {
			settings: Settings{
				InterfaceName: "$H1T",
			},
			errWrapped: ErrInterfaceNameInvalid,
			errMessage: "invalid interface name: $H1T",
		},
		"empty private key": {
			settings: Settings{
				InterfaceName: "wg0",
			},
			errWrapped: ErrPrivateKeyMissing,
			errMessage: "private key is missing",
		},
		"bad private key": {
			settings: Settings{
				InterfaceName: "wg0",
				PrivateKey:    "bad key",
			},
			errWrapped: ErrPrivateKeyInvalid,
			errMessage: "cannot parse private key",
		},
		"empty public key": {
			settings: Settings{
				InterfaceName: "wg0",
				PrivateKey:    validKey1,
			},
			errWrapped: ErrPublicKeyMissing,
			errMessage: "public key is missing",
		},
		"bad public key": {
			settings: Settings{
				InterfaceName: "wg0",
				PrivateKey:    validKey1,
				PublicKey:     "bad key",
			},
			errWrapped: ErrPublicKeyInvalid,
			errMessage: "cannot parse public key: bad key",
		},
		"bad preshared key": {
			settings: Settings{
				InterfaceName: "wg0",
				PrivateKey:    validKey1,
				PublicKey:     validKey2,
				PreSharedKey:  "bad key",
			},
			errWrapped: ErrPreSharedKeyInvalid,
			errMessage: "cannot parse pre-shared key",
		},
		"invalid endpoint address": {
			settings: Settings{
				InterfaceName: "wg0",
				PrivateKey:    validKey1,
				PublicKey:     validKey2,
			},
			errWrapped: ErrEndpointAddrMissing,
			errMessage: "endpoint address is missing",
		},
		"zero endpoint port": {
			settings: Settings{
				InterfaceName: "wg0",
				PrivateKey:    validKey1,
				PublicKey:     validKey2,
				Endpoint:      netip.AddrPortFrom(netip.AddrFrom4([4]byte{1, 2, 3, 4}), 0),
			},
			errWrapped: ErrEndpointPortMissing,
			errMessage: "endpoint port is missing",
		},
		"no address": {
			settings: Settings{
				InterfaceName: "wg0",
				PrivateKey:    validKey1,
				PublicKey:     validKey2,
				Endpoint:      netip.AddrPortFrom(netip.AddrFrom4([4]byte{1, 2, 3, 4}), 51820),
			},
			errWrapped: ErrAddressMissing,
			errMessage: "interface address is missing",
		},
		"invalid address": {
			settings: Settings{
				InterfaceName: "wg0",
				PrivateKey:    validKey1,
				PublicKey:     validKey2,
				Endpoint:      netip.AddrPortFrom(netip.AddrFrom4([4]byte{1, 2, 3, 4}), 51820),
				Addresses:     []netip.Prefix{{}},
			},
			errWrapped: ErrAddressNotValid,
			errMessage: "interface address is not valid: for address 1 of 1",
		},

		"no allowed IP": {
			settings: Settings{
				InterfaceName: "wg0",
				PrivateKey:    validKey1,
				PublicKey:     validKey2,
				Endpoint:      netip.AddrPortFrom(netip.AddrFrom4([4]byte{1, 2, 3, 4}), 51820),
				Addresses: []netip.Prefix{
					netip.PrefixFrom(netip.AddrFrom4([4]byte{5, 6, 7, 8}), 24),
				},
			},
			errWrapped: ErrAllowedIPsMissing,
			errMessage: "allowed IPs are missing",
		},
		"invalid allowed IP": {
			settings: Settings{
				InterfaceName: "wg0",
				PrivateKey:    validKey1,
				PublicKey:     validKey2,
				Endpoint:      netip.AddrPortFrom(netip.AddrFrom4([4]byte{1, 2, 3, 4}), 51820),
				Addresses: []netip.Prefix{
					netip.PrefixFrom(netip.AddrFrom4([4]byte{5, 6, 7, 8}), 24),
				},
				AllowedIPs: []netip.Prefix{{}},
			},
			errWrapped: ErrAllowedIPNotValid,
			errMessage: "allowed IP is not valid: for allowed IP 1 of 1",
		},
		"ipv6 allowed IP": {
			settings: Settings{
				InterfaceName: "wg0",
				PrivateKey:    validKey1,
				PublicKey:     validKey2,
				Endpoint:      netip.AddrPortFrom(netip.AddrFrom4([4]byte{1, 2, 3, 4}), 51820),
				Addresses: []netip.Prefix{
					netip.PrefixFrom(netip.AddrFrom4([4]byte{5, 6, 7, 8}), 24),
				},
				AllowedIPs: []netip.Prefix{
					allIPv6(),
				},
				IPv6: ptrTo(false),
			},
			errWrapped: ErrAllowedIPv6NotSupported,
			errMessage: "allowed IPv6 address not supported: for allowed IP ::/0",
		},
		"zero firewall mark": {
			settings: Settings{
				InterfaceName: "wg0",
				PrivateKey:    validKey1,
				PublicKey:     validKey2,
				Endpoint:      netip.AddrPortFrom(netip.AddrFrom4([4]byte{1, 2, 3, 4}), 51820),
				AllowedIPs:    []netip.Prefix{allIPv4()},
				Addresses: []netip.Prefix{
					netip.PrefixFrom(netip.AddrFrom4([4]byte{1, 2, 3, 4}), 24),
				},
			},
			errWrapped: ErrFirewallMarkMissing,
			errMessage: "firewall mark is missing",
		},
		"missing_MTU": {
			settings: Settings{
				InterfaceName: "wg0",
				PrivateKey:    validKey1,
				PublicKey:     validKey2,
				Endpoint:      netip.AddrPortFrom(netip.AddrFrom4([4]byte{1, 2, 3, 4}), 51820),
				AllowedIPs:    []netip.Prefix{allIPv4()},
				Addresses: []netip.Prefix{
					netip.PrefixFrom(netip.AddrFrom4([4]byte{1, 2, 3, 4}), 24),
				},
				FirewallMark: 999,
			},
			errWrapped: ErrMTUMissing,
			errMessage: "MTU is missing",
		},
		"invalid implementation": {
			settings: Settings{
				InterfaceName: "wg0",
				PrivateKey:    validKey1,
				PublicKey:     validKey2,
				Endpoint:      netip.AddrPortFrom(netip.AddrFrom4([4]byte{1, 2, 3, 4}), 51820),
				AllowedIPs:    []netip.Prefix{allIPv4()},
				Addresses: []netip.Prefix{
					netip.PrefixFrom(netip.AddrFrom4([4]byte{1, 2, 3, 4}), 24),
				},
				FirewallMark:   999,
				MTU:            1420,
				Implementation: "x",
			},
			errWrapped: ErrImplementationInvalid,
			errMessage: "invalid implementation: x",
		},
		"all valid": {
			settings: Settings{
				InterfaceName: "wg0",
				PrivateKey:    validKey1,
				PublicKey:     validKey2,
				Endpoint:      netip.AddrPortFrom(netip.AddrFrom4([4]byte{1, 2, 3, 4}), 51820),
				AllowedIPs: []netip.Prefix{
					allIPv6(),
				},
				Addresses: []netip.Prefix{
					netip.PrefixFrom(netip.AddrFrom4([4]byte{1, 2, 3, 4}), 24),
				},
				FirewallMark:   999,
				MTU:            1420,
				IPv6:           ptrTo(true),
				Implementation: "userspace",
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := testCase.settings.Check()

			assert.ErrorIs(t, err, testCase.errWrapped)
			if testCase.errWrapped != nil {
				assert.EqualError(t, err, testCase.errMessage)
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
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			lines := testCase.settings.ToLines(testCase.lineSettings)

			assert.Equal(t, testCase.lines, lines)
		})
	}
}
