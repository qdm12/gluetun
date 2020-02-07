package pia

import (
	"io/ioutil"
	"net"
	"strings"
	"testing"

	loggingMocks "github.com/qdm12/golibs/logging/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ModifyLines(t *testing.T) {
	t.Parallel()
	original, err := ioutil.ReadFile("testdata/ovpn.golden")
	require.NoError(t, err)
	originalLines := strings.Split(string(original), "\n")
	expected, err := ioutil.ReadFile("testdata/ovpn.modified.golden")
	require.NoError(t, err)
	expectedLines := strings.Split(string(expected), "\n")

	var port uint16 = 3000
	IPs := []net.IP{net.IP{100, 10, 10, 10}, net.IP{100, 20, 20, 20}}
	logger := &loggingMocks.Logger{}
	logger.On("Info", "%s: adapting openvpn configuration for server IP addresses and port %d", logPrefix, port).Once()
	c := &configurator{logger: logger}
	modifiedLines := c.ModifyLines(originalLines, IPs, port)
	assert.Equal(t, expectedLines, modifiedLines)
	logger.AssertExpectations(t)
}
