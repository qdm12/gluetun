package ivpn

import (
	"testing"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/stretchr/testify/assert"
)

func Test_Stringify(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		servers []models.IvpnServer
		s       string
	}{
		"no server": {
			s: `func IvpnServers() []models.IvpnServer {
	return []models.IvpnServer{
	}
}`,
		},
		"multiple servers": {
			servers: []models.IvpnServer{
				{Country: "A"},
				{Country: "B"},
			},
			s: `func IvpnServers() []models.IvpnServer {
	return []models.IvpnServer{
		{Country: "A", City: "", Hostname: "", IPs: []net.IP{}},
		{Country: "B", City: "", Hostname: "", IPs: []net.IP{}},
	}
}`,
		},
	}
	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			s := Stringify(testCase.servers)
			assert.Equal(t, testCase.s, s)
		})
	}
}
