package pia

import (
	"fmt"
	"io/ioutil"
	"net"
	"strings"
	"testing"

	loggingMocks "github.com/qdm12/golibs/logging/mocks"
	"github.com/qdm12/golibs/verification"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ParseConfig(t *testing.T) {
	t.Parallel()
	original, err := ioutil.ReadFile("testdata/ovpn.golden")
	require.NoError(t, err)
	exampleLines := strings.Split(string(original), "\n")
	tests := map[string]struct {
		lines       []string
		lookupIPs   []net.IP
		lookupIPErr error
		IPs         []net.IP
		port        uint16
		device      models.VPNDevice
		err         error
	}{
		"no data": {
			err: fmt.Errorf("remote line not found in Openvpn configuration"),
		},
		"bad remote line": {
			lines: []string{"remote field2"},
			err:   fmt.Errorf("line \"remote field2\" misses information"),
		},
		"bad remote port": {
			lines: []string{"remote field2 port"},
			err:   fmt.Errorf("line \"remote field2 port\" has an invalid port: port \"port\" is not a valid integer"),
		},
		"lookupIP error": {
			lines:       []string{"remote host 1000"},
			lookupIPErr: fmt.Errorf("lookup error"),
			err:         fmt.Errorf("lookup error"),
		},
		"missing dev line": {
			lines: []string{"remote host 1994"},
			err:   fmt.Errorf("device line not found in Openvpn configuration"),
		},
		"bad dev line": {
			lines: []string{"dev   field2 field3"},
			err:   fmt.Errorf("line \"dev   field2 field3\" misses information"),
		},
		"bad device": {
			lines: []string{"dev xx"},
			err:   fmt.Errorf("device \"xx0\" is not valid"),
		},
		"valid lines": {
			lines:  []string{"remote host 1194", "dev tap", "blabla"},
			port:   1194,
			device: constants.TAP,
		},
		"real data": {
			lines:     exampleLines,
			lookupIPs: []net.IP{{100, 100, 100, 100}, {100, 100, 200, 200}},
			IPs:       []net.IP{{100, 100, 100, 100}, {100, 100, 200, 200}},
			port:      1198,
			device:    constants.TUN,
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			logger := &loggingMocks.Logger{}
			logger.On("Info", "%s: parsing openvpn configuration", logPrefix).Once()
			if tc.err == nil {
				logger.On("Info", "%s: Found %d PIA server IP addresses, port %d and device %s", logPrefix, len(tc.IPs), tc.port, tc.device).Once()
			}
			lookupIP := func(host string) ([]net.IP, error) {
				return tc.lookupIPs, tc.lookupIPErr
			}
			c := &configurator{logger: logger, verifyPort: verification.NewVerifier().VerifyPort, lookupIP: lookupIP}
			IPs, port, device, err := c.ParseConfig(tc.lines)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.IPs, IPs)
			assert.Equal(t, tc.port, port)
			assert.Equal(t, tc.device, device)
			logger.AssertExpectations(t)
		})
	}
}
