package privateinternetaccess

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_newHTTPClient(t *testing.T) {
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

	piaClient, err := newHTTPClient(serverName)

	require.NoError(t, err)

	// Verify pia transport TLS config is set
	piaTransport := piaClient.Transport.(*http.Transport)
	rootCAsSubjects := piaTransport.TLSClientConfig.RootCAs.Subjects()
	assert.Equal(t, expectedRootCAsSubjects, rootCAsSubjects)
	piaTransport.TLSClientConfig.RootCAs = nil
	assert.Equal(t, expectedPIATransportTLSConfig, piaTransport.TLSClientConfig)
}

func Test_unpackPayload(t *testing.T) {
	t.Parallel()

	const exampleToken = "token"
	const examplePort = 2000
	exampleExpiration := time.Unix(1000, 0).UTC()

	testCases := map[string]struct {
		payload    string
		port       uint16
		token      string
		expiration time.Time
		err        error
	}{
		"valid payload": {
			payload:    makePIAPayload(t, exampleToken, examplePort, exampleExpiration),
			port:       examplePort,
			token:      exampleToken,
			expiration: exampleExpiration,
			err:        nil,
		},
		"invalid base64 payload": {
			payload: "invalid",
			err:     errors.New("illegal base64 data at input byte 4: for payload: invalid"),
		},
		"invalid json payload": {
			payload: base64.StdEncoding.EncodeToString([]byte{1}),
			err:     errors.New("invalid character '\\x01' looking for beginning of value: for data: \x01"),
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			port, token, expiration, err := unpackPayload(testCase.payload)

			if testCase.err != nil {
				require.Error(t, err)
				assert.Equal(t, testCase.err.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, testCase.port, port)
			assert.Equal(t, testCase.token, token)
			assert.Equal(t, testCase.expiration, expiration)
		})
	}
}

func makePIAPayload(t *testing.T, token string, port uint16, expiration time.Time) (payload string) {
	t.Helper()

	data := piaPayload{
		Token:      token,
		Port:       port,
		Expiration: expiration,
	}

	b, err := json.Marshal(data)
	require.NoError(t, err)

	return base64.StdEncoding.EncodeToString(b)
}

func Test_replaceInString(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		s             string
		substitutions map[string]string
		result        string
	}{
		"empty": {},
		"multiple replacements": {
			s: "https://test.com/username/password/",
			substitutions: map[string]string{
				"username": "xxx",
				"password": "yyy",
			},
			result: "https://test.com/xxx/yyy/",
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			result := replaceInString(testCase.s, testCase.substitutions)
			assert.Equal(t, testCase.result, result)
		})
	}
}
