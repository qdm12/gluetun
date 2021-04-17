package provider

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"net/http"
	"testing"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_newPIAHTTPClient(t *testing.T) {
	t.Parallel()

	const serverName = "testserver"

	certificateBytes, err := base64.StdEncoding.DecodeString(constants.PIACertificateStrong)
	require.NoError(t, err)
	certificate, err := x509.ParseCertificate(certificateBytes)
	require.NoError(t, err)
	rootCAs := x509.NewCertPool()
	rootCAs.AddCert(certificate)
	expectedRootCAsSubjects := rootCAs.Subjects()

	expectedPIATransportTLSConfig := &tls.Config{
		// Can't directly compare RootCAs because of private fields
		RootCAs:    nil,
		MinVersion: tls.VersionTLS12,
		ServerName: serverName,
	}

	piaClient, err := newPIAHTTPClient(serverName)

	require.NoError(t, err)

	// Verify pia transport TLS config is set
	piaTransport := piaClient.Transport.(*http.Transport)
	rootCAsSubjects := piaTransport.TLSClientConfig.RootCAs.Subjects()
	assert.Equal(t, expectedRootCAsSubjects, rootCAsSubjects)
	piaTransport.TLSClientConfig.RootCAs = nil
	assert.Equal(t, expectedPIATransportTLSConfig, piaTransport.TLSClientConfig)
}
