package settings

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Unbound_JSON(t *testing.T) {
	t.Parallel()

	settings := Unbound{
		Providers: []string{"cloudflare"},
		Caching:   boolPtr(true),
		IPv6:      boolPtr(false),
	}

	b, err := json.Marshal(settings)
	require.NoError(t, err)

	const expected = `{"providers":["cloudflare"],"caching":true,"ipv6":false}`

	assert.Equal(t, expected, string(b))

	var resultSettings Unbound
	err = json.Unmarshal(b, &resultSettings)
	require.NoError(t, err)

	assert.Equal(t, settings, resultSettings)
}
