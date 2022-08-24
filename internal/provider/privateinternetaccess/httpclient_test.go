package privateinternetaccess

import (
	"crypto/tls"
	"crypto/x509/pkix"
	"encoding/asn1"
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

	piaClient, err := newHTTPClient(serverName)
	require.NoError(t, err)

	// Verify pia transport TLS config is set
	piaTransport, ok := piaClient.Transport.(*http.Transport)
	require.True(t, ok)

	subjects := piaTransport.TLSClientConfig.RootCAs.Subjects() //nolint:staticcheck
	assert.NotEmpty(t, subjects)
	piaCertFound := false
	for _, subject := range subjects {
		var rdnSequence pkix.RDNSequence
		_, err := asn1.Unmarshal(subject, &rdnSequence)
		require.NoError(t, err)
		var name pkix.Name
		name.FillFromRDNSequence(&rdnSequence)
		if name.CommonName == "Private Internet Access" {
			piaCertFound = true
			break
		}
	}
	assert.True(t, piaCertFound)

	piaTransport.TLSClientConfig.RootCAs = nil
	assert.Equal(t, expectedPIATransportTLSConfig, piaTransport.TLSClientConfig)
}
