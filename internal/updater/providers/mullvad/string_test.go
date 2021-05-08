package mullvad

import (
	"net"
	"strings"
	"testing"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/stretchr/testify/assert"
)

func Test_Stringify(t *testing.T) {
	servers := []models.MullvadServer{{
		Country:  "webland",
		City:     "webcity",
		Hostname: "hostname",
		ISP:      "not nsa",
		Owned:    true,
		IPs:      []net.IP{{1, 1, 1, 1}},
		IPsV6:    []net.IP{{1, 1, 1, 1}},
	}}
	//nolint:lll
	expected := `
func MullvadServers() []models.MullvadServer {
	return []models.MullvadServer{
		{Country: "webland", City: "webcity", Hostname: "hostname", ISP: "not nsa", Owned: true, IPs: []net.IP{{1, 1, 1, 1}}, IPsV6: []net.IP{{1, 1, 1, 1}}},
	}
}
`
	expected = strings.TrimPrefix(strings.TrimSuffix(expected, "\n"), "\n")
	s := Stringify(servers)
	assert.Equal(t, expected, s)
}
