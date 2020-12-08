package settings

import (
	"encoding/json"
	"testing"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_OpenVPN_JSON(t *testing.T) {
	t.Parallel()
	in := OpenVPN{
		Root: true,
		Provider: models.ProviderSettings{
			Name: "name",
		},
	}
	data, err := json.Marshal(in)
	require.NoError(t, err)
	//nolint:lll
	assert.Equal(t, `{"user":"","password":"","verbosity":0,"runAsRoot":true,"cipher":"","auth":"","provider":{"name":"name","serverSelection":{"networkProtocol":"","regions":null,"group":"","countries":null,"cities":null,"hostnames":null,"isps":null,"owned":false,"customPort":0,"numbers":null,"encryptionPreset":""},"extraConfig":{"encryptionPreset":"","openvpnIPv6":false},"portForwarding":{"enabled":false,"filepath":""}}}`, string(data))
	var out OpenVPN
	err = json.Unmarshal(data, &out)
	require.NoError(t, err)
	assert.Equal(t, in, out)
}
