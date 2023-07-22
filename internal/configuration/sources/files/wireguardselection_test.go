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

func uint16Ptr(n uint16) *uint16 { return &n }

func Test_Source_readWireguardSelection(t *testing.T) {
	t.Parallel()

	t.Run("fail reading from file", func(t *testing.T) {
		t.Parallel()

		dirPath := t.TempDir()
		source := &Source{
			wireguardConfigPath: dirPath,
		}
		wireguard, err := source.readWireguardSelection()
		assert.Equal(t, settings.WireguardSelection{}, wireguard)
		assert.Error(t, err)
		assert.Regexp(t, `reading file: read .+: is a directory`, err.Error())
	})

	t.Run("no file", func(t *testing.T) {
		t.Parallel()

		noFile := filepath.Join(t.TempDir(), "doesnotexist")
		source := &Source{
			wireguardConfigPath: noFile,
		}
		wireguard, err := source.readWireguardSelection()
		assert.Equal(t, settings.WireguardSelection{}, wireguard)
		assert.NoError(t, err)
	})

	testCases := map[string]struct {
		fileContent string
		selection   settings.WireguardSelection
		errMessage  string
	}{
		"ini load error": {
			fileContent: "invalid",
			errMessage:  "loading ini from reader: key-value delimiter not found: invalid",
		},
		"empty file": {},
		"peer section parsing error": {
			fileContent: `
[Peer]
PublicKey = x
`,
			errMessage: "parsing peer section: parsing PublicKey: " +
				"x: wgtypes: failed to parse base64-encoded key: " +
				"illegal base64 data at input byte 0",
		},
		"success": {
			fileContent: `
[Peer]
PublicKey = QOlCgyA/Sn/c/+YNTIEohrjm8IZV+OZ2AUFIoX20sk8=
Endpoint = 1.2.3.4:51820
`,
			selection: settings.WireguardSelection{
				PublicKey:    "QOlCgyA/Sn/c/+YNTIEohrjm8IZV+OZ2AUFIoX20sk8=",
				EndpointIP:   netip.AddrFrom4([4]byte{1, 2, 3, 4}),
				EndpointPort: uint16Ptr(51820),
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

			wireguard, err := source.readWireguardSelection()

			assert.Equal(t, testCase.selection, wireguard)
			if testCase.errMessage != "" {
				assert.EqualError(t, err, testCase.errMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_parseWireguardPeerSection(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		iniData    string
		selection  settings.WireguardSelection
		errMessage string
	}{
		"public key error": {
			iniData: `[Peer]
PublicKey = x`,
			errMessage: "parsing PublicKey: x: " +
				"wgtypes: failed to parse base64-encoded key: " +
				"illegal base64 data at input byte 0",
		},
		"public key set": {
			iniData: `[Peer]
PublicKey = QOlCgyA/Sn/c/+YNTIEohrjm8IZV+OZ2AUFIoX20sk8=`,
			selection: settings.WireguardSelection{
				PublicKey: "QOlCgyA/Sn/c/+YNTIEohrjm8IZV+OZ2AUFIoX20sk8=",
			},
		},
		"missing port in endpoint": {
			iniData: `[Peer]
Endpoint = x`,
			errMessage: "splitting endpoint: address x: missing port in address",
		},
		"endpoint host is not IP": {
			iniData: `[Peer]
Endpoint = website.com:51820`,
			errMessage: "endpoint host is not an IP: ParseAddr(\"website.com\"): unexpected character (at \"website.com\")",
		},
		"endpoint port is not valid": {
			iniData: `[Peer]
Endpoint = 1.2.3.4:518299`,
			errMessage: "port from Endpoint key: port cannot be higher than 65535: 518299",
		},
		"valid endpoint": {
			iniData: `[Peer]
Endpoint = 1.2.3.4:51820`,
			selection: settings.WireguardSelection{
				EndpointIP:   netip.AddrFrom4([4]byte{1, 2, 3, 4}),
				EndpointPort: uint16Ptr(51820),
			},
		},
		"all set": {
			iniData: `[Peer]
PublicKey = QOlCgyA/Sn/c/+YNTIEohrjm8IZV+OZ2AUFIoX20sk8=
Endpoint = 1.2.3.4:51820`,
			selection: settings.WireguardSelection{
				PublicKey:    "QOlCgyA/Sn/c/+YNTIEohrjm8IZV+OZ2AUFIoX20sk8=",
				EndpointIP:   netip.AddrFrom4([4]byte{1, 2, 3, 4}),
				EndpointPort: uint16Ptr(51820),
			},
		},
	}

	for testName, testCase := range testCases {
		testCase := testCase
		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			iniFile, err := ini.Load([]byte(testCase.iniData))
			require.NoError(t, err)
			iniSection, err := iniFile.GetSection("Peer")
			require.NoError(t, err)

			var selection settings.WireguardSelection
			err = parseWireguardPeerSection(iniSection, &selection)

			assert.Equal(t, testCase.selection, selection)
			if testCase.errMessage != "" {
				assert.EqualError(t, err, testCase.errMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
