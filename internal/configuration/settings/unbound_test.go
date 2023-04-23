package settings

import (
	"encoding/json"
	"net/netip"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Unbound_JSON(t *testing.T) {
	t.Parallel()

	settings := Unbound{
		Providers:             []string{"cloudflare"},
		Caching:               boolPtr(true),
		IPv6:                  boolPtr(false),
		VerbosityLevel:        uint8Ptr(1),
		VerbosityDetailsLevel: nil,
		ValidationLogLevel:    uint8Ptr(0),
		Username:              "user",
		Allowed: []netip.Prefix{
			netip.PrefixFrom(netip.AddrFrom4([4]byte{}), 0),
			netip.PrefixFrom(netip.AddrFrom16([16]byte{}), 0),
		},
	}

	b, err := json.Marshal(settings)
	require.NoError(t, err)

	const expected = `{"Providers":["cloudflare"],"Caching":true,"IPv6":false,` +
		`"VerbosityLevel":1,"VerbosityDetailsLevel":null,"ValidationLogLevel":0,` +
		`"Username":"user","Allowed":["0.0.0.0/0","::/0"]}`

	assert.Equal(t, expected, string(b))

	var resultSettings Unbound
	err = json.Unmarshal(b, &resultSettings)
	require.NoError(t, err)

	assert.Equal(t, settings, resultSettings)
}
