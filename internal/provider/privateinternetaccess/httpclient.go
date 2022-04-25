package privateinternetaccess

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

func newHTTPClient(serverName string) (client *http.Client) {
	//nolint:gomnd
	return &http.Client{
		Transport: &http.Transport{
			// Settings taken from http.DefaultTransport
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
				ServerName: serverName,
			},
		},
		Timeout: 30 * time.Second,
	}
}
