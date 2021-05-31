package ipvanish

import (
	"testing"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/stretchr/testify/assert"
)

func Test_Stringify(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		servers []models.IpvanishServer
		s       string
	}{
		"no server": {
			s: `func IpvanishServers() []models.IpvanishServer {
	return []models.IpvanishServer{
	}
}`,
		},
		"multiple servers": {
			servers: []models.IpvanishServer{
				{Country: "A"},
				{Country: "B"},
			},
			s: `func IpvanishServers() []models.IpvanishServer {
	return []models.IpvanishServer{
		{Country: "A", City: "", Hostname: "", TCP: false, UDP: false, IPs: []net.IP{}},
		{Country: "B", City: "", Hostname: "", TCP: false, UDP: false, IPs: []net.IP{}},
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
