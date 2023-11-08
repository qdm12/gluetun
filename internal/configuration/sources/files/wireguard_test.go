package files

import (
	"net/netip"
	"os"
	"path/filepath"
	"testing"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/ini.v1"
)

func Test_Source_readWireguard(t *testing.T) {
	t.Parallel()

	t.Run("fail reading from file", func(t *testing.T) {
		t.Parallel()

		dirPath := t.TempDir()
		source := &Source{
			wireguardConfigPath: dirPath,
		}
		wireguard, err := source.readWireguard()
		assert.Equal(t, settings.Wireguard{}, wireguard)
		assert.Error(t, err)
		assert.Regexp(t, `reading file: read .+: is a directory`, err.Error())
	})

	t.Run("no file", func(t *testing.T) {
		t.Parallel()

		noFile := filepath.Join(t.TempDir(), "doesnotexist")
		source := &Source{
			wireguardConfigPath: noFile,
		}
		wireguard, err := source.readWireguard()
		assert.Equal(t, settings.Wireguard{}, wireguard)
		assert.NoError(t, err)
	})

	testCases := map[string]struct {
		fileContent string
		wireguard   settings.Wireguard
		errMessage  string
	}{
		"ini load error": {
			fileContent: "invalid",
			errMessage:  "loading ini from reader: key-value delimiter not found: invalid",
		},
		"empty file": {},
		"interface section parsing error": {
			fileContent: `
[Interface]
PrivateKey = x
`,
			errMessage: "parsing interface section: parsing PrivateKey: " +
				"x: wgtypes: failed to parse base64-encoded key: " +
				"illegal base64 data at input byte 0",
		},
		"success": {
			fileContent: `
[Interface]
PrivateKey = QOlCgyA/Sn/c/+YNTIEohrjm8IZV+OZ2AUFIoX20sk8=
PresharedKey = YJ680VN+dGrdsWNjSFqZ6vvwuiNhbq502ZL3G7Q3o3g=
Address = 10.38.22.35/32
DNS = 193.138.218.74

[Peer]
`,
			wireguard: settings.Wireguard{
				PrivateKey:   ptrTo("QOlCgyA/Sn/c/+YNTIEohrjm8IZV+OZ2AUFIoX20sk8="),
				PreSharedKey: ptrTo("YJ680VN+dGrdsWNjSFqZ6vvwuiNhbq502ZL3G7Q3o3g="),
				Addresses: []netip.Prefix{
					netip.PrefixFrom(netip.AddrFrom4([4]byte{10, 38, 22, 35}), 32),
				},
			},
		},
	}

	for testName, testCase := range testCases {
		testCase := testCase
		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			configFile := filepath.Join(t.TempDir(), "wg.conf")
			err := os.WriteFile(configFile, []byte(testCase.fileContent), 0600)
			require.NoError(t, err)

			source := &Source{
				wireguardConfigPath: configFile,
			}

			wireguard, err := source.readWireguard()

			assert.Equal(t, testCase.wireguard, wireguard)
			if testCase.errMessage != "" {
				assert.EqualError(t, err, testCase.errMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_parseWireguardInterfaceSection(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		iniData    string
		wireguard  settings.Wireguard
		errMessage string
	}{
		"private key error": {
			iniData: `[Interface]
PrivateKey = x`,
			errMessage: "parsing PrivateKey: x: " +
				"wgtypes: failed to parse base64-encoded key: " +
				"illegal base64 data at input byte 0",
		},
		"pre shared key error": {
			iniData: `[Interface]
PresharedKey = x
`,
			errMessage: "parsing PresharedKey: x: " +
				"wgtypes: failed to parse base64-encoded key: " +
				"illegal base64 data at input byte 0",
		},
		"address error": {
			iniData: `[Interface]
Address = x
`,
			errMessage: "parsing address: netip.ParsePrefix(\"x/32\"): ParseAddr(\"x\"): unable to parse IP",
		},
		"success": {
			iniData: `
[Interface]
PrivateKey = QOlCgyA/Sn/c/+YNTIEohrjm8IZV+OZ2AUFIoX20sk8=
PresharedKey = YJ680VN+dGrdsWNjSFqZ6vvwuiNhbq502ZL3G7Q3o3g=
Address = 10.38.22.35/32
`,
			wireguard: settings.Wireguard{
				PrivateKey:   ptrTo("QOlCgyA/Sn/c/+YNTIEohrjm8IZV+OZ2AUFIoX20sk8="),
				PreSharedKey: ptrTo("YJ680VN+dGrdsWNjSFqZ6vvwuiNhbq502ZL3G7Q3o3g="),
				Addresses: []netip.Prefix{
					netip.PrefixFrom(netip.AddrFrom4([4]byte{10, 38, 22, 35}), 32),
				},
			},
		},
	}

	for testName, testCase := range testCases {
		testCase := testCase
		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			iniFile, err := ini.Load([]byte(testCase.iniData))
			require.NoError(t, err)
			iniSection, err := iniFile.GetSection("Interface")
			require.NoError(t, err)

			var wireguard settings.Wireguard
			err = parseWireguardInterfaceSection(iniSection, &wireguard)

			assert.Equal(t, testCase.wireguard, wireguard)
			if testCase.errMessage != "" {
				assert.EqualError(t, err, testCase.errMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_parseINIWireguardKey(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		fileContent string
		keyName     string
		key         *string
		errMessage  string
	}{
		"key does not exist": {
			fileContent: `[Interface]`,
			keyName:     "PrivateKey",
		},
		"bad Wireguard key": {
			fileContent: `[Interface]
PrivateKey = x`,
			keyName: "PrivateKey",
			errMessage: "parsing PrivateKey: x: " +
				"wgtypes: failed to parse base64-encoded key: " +
				"illegal base64 data at input byte 0",
		},
		"success": {
			fileContent: `[Interface]
PrivateKey = QOlCgyA/Sn/c/+YNTIEohrjm8IZV+OZ2AUFIoX20sk8=`,
			keyName: "PrivateKey",
			key:     ptrTo("QOlCgyA/Sn/c/+YNTIEohrjm8IZV+OZ2AUFIoX20sk8="),
		},
	}

	for testName, testCase := range testCases {
		testCase := testCase
		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			iniFile, err := ini.Load([]byte(testCase.fileContent))
			require.NoError(t, err)
			iniSection, err := iniFile.GetSection("Interface")
			require.NoError(t, err)

			key, err := parseINIWireguardKey(iniSection, testCase.keyName)

			assert.Equal(t, testCase.key, key)
			if testCase.errMessage != "" {
				assert.EqualError(t, err, testCase.errMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_parseINIWireguardAddress(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		fileContent string
		addresses   []netip.Prefix
		errMessage  string
	}{
		"key does not exist": {
			fileContent: `[Interface]`,
		},
		"bad address": {
			fileContent: `[Interface]
Address = x`,
			errMessage: "parsing address: netip.ParsePrefix(\"x/32\"): ParseAddr(\"x\"): unable to parse IP",
		},
		"success": {
			fileContent: `[Interface]
Address = 1.2.3.4/32, 5.6.7.8/32`,
			addresses: []netip.Prefix{
				netip.PrefixFrom(netip.AddrFrom4([4]byte{1, 2, 3, 4}), 32),
				netip.PrefixFrom(netip.AddrFrom4([4]byte{5, 6, 7, 8}), 32),
			},
		},
	}

	for testName, testCase := range testCases {
		testCase := testCase
		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			iniFile, err := ini.Load([]byte(testCase.fileContent))
			require.NoError(t, err)
			iniSection, err := iniFile.GetSection("Interface")
			require.NoError(t, err)

			addresses, err := parseINIWireguardAddress(iniSection)

			assert.Equal(t, testCase.addresses, addresses)
			if testCase.errMessage != "" {
				assert.EqualError(t, err, testCase.errMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
