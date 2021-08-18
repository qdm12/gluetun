package configuration

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_OpenVPN_JSON(t *testing.T) {
	t.Parallel()
	in := OpenVPN{
		Root:  true,
		Flags: []string{},
	}
	data, err := json.MarshalIndent(in, "", "  ")
	require.NoError(t, err)
	assert.Equal(t, `{
  "user": "",
  "password": "",
  "verbosity": 0,
  "flags": [],
  "mssfix": 0,
  "run_as_root": true,
  "cipher": "",
  "auth": "",
  "custom_config": "",
  "version": "",
  "encryption_preset": "",
  "ipv6": false,
  "procuser": ""
}`, string(data))
	var out OpenVPN
	err = json.Unmarshal(data, &out)
	require.NoError(t, err)
	assert.Equal(t, in, out)
}
