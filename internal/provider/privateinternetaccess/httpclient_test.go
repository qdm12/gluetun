package privateinternetaccess

import (
	"crypto/tls"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_newHTTPClient(t *testing.T) {
	t.Parallel()

	const serverName = "testserver"

	expectedPIATransportTLSConfig := &tls.Config{
		// Can't directly compare RootCAs because of private fields
		RootCAs:    nil,
		MinVersion: tls.VersionTLS12,
		ServerName: serverName,
	}

	piaClient := newHTTPClient(serverName)

	// Verify pia transport TLS config is set
	piaTransport, ok := piaClient.Transport.(*http.Transport)
	require.True(t, ok)
	piaTransport.TLSClientConfig.RootCAs = nil
	assert.Equal(t, expectedPIATransportTLSConfig, piaTransport.TLSClientConfig)
}
