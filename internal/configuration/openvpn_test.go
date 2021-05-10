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
		Root: true,
		Provider: Provider{
			Name: "name",
		},
	}
	data, err := json.MarshalIndent(in, "", "  ")
	require.NoError(t, err)
	assert.Equal(t, `{
  "user": "",
  "password": "",
  "verbosity": 0,
  "mssfix": 0,
  "run_as_root": true,
  "cipher": "",
  "auth": "",
  "provider": {
    "name": "name",
    "server_selection": {
      "tcp": false,
      "regions": null,
      "group": "",
      "countries": null,
      "cities": null,
      "hostnames": null,
      "names": null,
      "isps": null,
      "owned": false,
      "custom_port": 0,
      "numbers": null,
      "encryption_preset": ""
    },
    "extra_config": {
      "encryption_preset": "",
      "openvpn_ipv6": false
    },
    "port_forwarding": {
      "enabled": false,
      "filepath": ""
    }
  },
  "custom_config": ""
}`, string(data))
	var out OpenVPN
	err = json.Unmarshal(data, &out)
	require.NoError(t, err)
	assert.Equal(t, in, out)
}
