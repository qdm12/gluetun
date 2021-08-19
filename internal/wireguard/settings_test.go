package wireguard

import (
	"errors"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Settings_SetDefaults(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		original Settings
		expected Settings
	}{
		"empty settings": {
			expected: Settings{
				InterfaceName: "wg0",
				FirewallMark:  51820,
			},
		},
		"not empty settings": {
			original: Settings{
				InterfaceName: "wg1",
				FirewallMark:  999,
			},
			expected: Settings{
				InterfaceName: "wg1",
				FirewallMark:  999,
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
		"empty endpoint": {
			settings: Settings{
				InterfaceName: "wg0",
				PrivateKey:    validKey1,
				PublicKey:     validKey2,
			},
			err: ErrEndpointMissing,
		},
		"nil endpoint IP": {
			settings: Settings{
				InterfaceName: "wg0",
				PrivateKey:    validKey1,
				PublicKey:     validKey2,
				Endpoint:      &net.UDPAddr{},
			},
			err: ErrEndpointIPMissing,
		},
		"nil endpoint port": {
			settings: Settings{
				InterfaceName: "wg0",
				PrivateKey:    validKey1,
				PublicKey:     validKey2,
				Endpoint: &net.UDPAddr{
					IP: net.IPv4(1, 2, 3, 4),
				},
			},
			err: ErrEndpointPortMissing,
		},
		"no address": {
			settings: Settings{
				InterfaceName: "wg0",
				PrivateKey:    validKey1,
				PublicKey:     validKey2,
				Endpoint: &net.UDPAddr{
					IP:   net.IPv4(1, 2, 3, 4),
					Port: 51820,
				},
			},
			err: ErrAddressMissing,
		},
		"nil address": {
			settings: Settings{
				InterfaceName: "wg0",
				PrivateKey:    validKey1,
				PublicKey:     validKey2,
				Endpoint: &net.UDPAddr{
					IP:   net.IPv4(1, 2, 3, 4),
					Port: 51820,
				},
				Addresses: []*net.IPNet{nil},
			},
			err: errors.New("interface address is nil: for address 1 of 1"),
		},
		"nil address IP": {
			settings: Settings{
				InterfaceName: "wg0",
				PrivateKey:    validKey1,
				PublicKey:     validKey2,
				Endpoint: &net.UDPAddr{
					IP:   net.IPv4(1, 2, 3, 4),
					Port: 51820,
				},
				Addresses: []*net.IPNet{{}},
			},
			err: errors.New("interface address IP is missing: for address 1 of 1"),
		},
		"nil address mask": {
			settings: Settings{
				InterfaceName: "wg0",
				PrivateKey:    validKey1,
				PublicKey:     validKey2,
				Endpoint: &net.UDPAddr{
					IP:   net.IPv4(1, 2, 3, 4),
					Port: 51820,
				},
				Addresses: []*net.IPNet{{IP: net.IPv4(1, 2, 3, 4)}},
			},
			err: errors.New("interface address mask is missing: for address 1 of 1"),
		},
		"zero firewall mark": {
			settings: Settings{
				InterfaceName: "wg0",
				PrivateKey:    validKey1,
				PublicKey:     validKey2,
				Endpoint: &net.UDPAddr{
					IP:   net.IPv4(1, 2, 3, 4),
					Port: 51820,
				},
				Addresses: []*net.IPNet{{IP: net.IPv4(1, 2, 3, 4), Mask: net.CIDRMask(24, 32)}},
			},
			err: ErrFirewallMarkMissing,
		},
		"all valid": {
			settings: Settings{
				InterfaceName: "wg0",
				PrivateKey:    validKey1,
				PublicKey:     validKey2,
				Endpoint: &net.UDPAddr{
					IP:   net.IPv4(1, 2, 3, 4),
					Port: 51820,
				},
				Addresses:    []*net.IPNet{{IP: net.IPv4(1, 2, 3, 4), Mask: net.CIDRMask(24, 32)}},
				FirewallMark: 999,
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

func Test_Settings_Lines(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		settings     Settings
		lineSettings ToLinesSettings
		lines        []string
	}{
		"empty settings": {
			lines: []string{
				"├── Interface name: ",
				"├── Private key: not set",
				"├── Pre shared key: not set",
				"├── Endpoint: not set",
				"└── Addresses: not set",
			},
		},
		"settings all set": {
			settings: Settings{
				InterfaceName: "wg0",
				PrivateKey:    "private key",
				PublicKey:     "public key",
				PreSharedKey:  "pre-shared key",
				Endpoint: &net.UDPAddr{
					IP:   net.IPv4(1, 2, 3, 4),
					Port: 51820,
				},
				FirewallMark: 999,
				Addresses: []*net.IPNet{
					{IP: net.IPv4(1, 1, 1, 1), Mask: net.CIDRMask(24, 32)},
					{IP: net.IPv4(2, 2, 2, 2), Mask: net.CIDRMask(32, 32)},
				},
			},
			lines: []string{
				"├── Interface name: wg0",
				"├── Private key: set",
				"├── PublicKey: public key",
				"├── Pre shared key: set",
				"├── Endpoint: 1.2.3.4:51820",
				"├── Firewall mark: 999",
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
				Addresses: []*net.IPNet{
					{IP: net.IPv4(1, 1, 1, 1), Mask: net.CIDRMask(24, 32)},
					{IP: net.IPv4(2, 2, 2, 2), Mask: net.CIDRMask(32, 32)},
				},
			},
			lines: []string{
				"- Interface name: wg0",
				"- Private key: not set",
				"- Pre shared key: not set",
				"- Endpoint: not set",
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
