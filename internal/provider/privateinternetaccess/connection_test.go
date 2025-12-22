package privateinternetaccess

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_getConnectionDefaults(t *testing.T) {
	t.Parallel()

	const timeout = 5 * time.Second
	client := &http.Client{
		Timeout: timeout,
	}

	ctx := t.Context()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet,
		"https://serverlist.piaservers.net/vpninfo/servers/v6", nil)
	require.NoError(t, err)

	response, err := client.Do(request)
	require.NoError(t, err)
	defer response.Body.Close()

	require.Equal(t, http.StatusOK, response.StatusCode)

	b, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	// remove key/signature at the bottom
	i := bytes.IndexRune(b, '\n')
	b = b[:i]

	var data struct {
		Groups struct {
			OvpnUDP []struct {
				Ports []uint16 `json:"ports"`
			} `json:"ovpnudp"`
			OvpnTCP []struct {
				Ports []uint16 `json:"ports"`
			} `json:"ovpntcp"`
		} `json:"groups"`
	}
	err = json.Unmarshal(b, &data)
	require.NoError(t, err)

	defaults := getConnectionDefaults()

	require.Len(t, data.Groups.OvpnUDP, 1)
	require.Len(t, data.Groups.OvpnTCP, 1)
	assert.Contains(t, data.Groups.OvpnUDP[0].Ports, defaults.OpenVPNUDPPort)
	assert.Contains(t, data.Groups.OvpnTCP[0].Ports, defaults.OpenVPNTCPPort)
}
