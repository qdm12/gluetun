package files

import (
	"net/netip"
	"testing"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Source_readWireguard_integration(t *testing.T) {
	t.Parallel()

	source := &Source{
		wireguardConfigPath: "./testdata/wg.conf",
	}

	wireguard, err := source.readWireguard()
	require.NoError(t, err)

	expectedWireguard := settings.Wireguard{
		PrivateKey:   ptrTo("QOlCgyA/Sn/c/+YNTIEohrjm8IZV+OZ2AUFIoX20sk8="),
		PreSharedKey: ptrTo("YJ680VN+dGrdsWNjSFqZ6vvwuiNhbq502ZL3G7Q3o3g="),
		Addresses: []netip.Prefix{
			netip.PrefixFrom(netip.AddrFrom4([4]byte{10, 38, 22, 35}), 32),
			netip.PrefixFrom(netip.AddrFrom16([16]byte{
				0xfa, 0x00, 0xdd, 0xdd, 0xbc, 0xcc, 0xbb, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, 0xe2, 0x22,
			}), 128),
		},
	}
	assert.Equal(t, expectedWireguard, wireguard)
}
