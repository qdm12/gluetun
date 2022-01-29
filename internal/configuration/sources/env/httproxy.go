package env

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/govalid/binary"
)

func (r *Reader) readHTTPProxy() (httpProxy settings.HTTPProxy, err error) {
	httpProxy.User = r.readHTTProxyUser()
	httpProxy.Password = r.readHTTProxyPassword()
	httpProxy.ListeningAddress = r.readHTTProxyListeningAddress()

	httpProxy.Enabled, err = r.readHTTProxyEnabled()
	if err != nil {
		return httpProxy, err
	}

	httpProxy.Stealth, err = envToBoolPtr("HTTPPROXY_STEALTH")
	if err != nil {
		return httpProxy, fmt.Errorf("environment variable HTTPPROXY_STEALTH: %w", err)
	}

	httpProxy.Log, err = r.readHTTProxyLog()
	if err != nil {
		return httpProxy, err
	}

	return httpProxy, nil
}

func (r *Reader) readHTTProxyUser() (user *string) {
	_, s := r.getEnvWithRetro("HTTPPROXY_USER", "PROXY_USER", "TINYPROXY_USER")
	if s != "" {
		return &s
	}
	return nil
}

func (r *Reader) readHTTProxyPassword() (user *string) {
	_, s := r.getEnvWithRetro("HTTPPROXY_PASSWORD", "PROXY_PASSWORD", "TINYPROXY_PASSWORD")
	if s != "" {
		return &s
	}
	return nil
}

func (r *Reader) readHTTProxyListeningAddress() (listeningAddress string) {
	key, value := r.getEnvWithRetro("HTTPPROXY_LISTENING_ADDRESS", "PROXY_PORT", "TINYPROXY_PORT", "HTTPPROXY_PORT")
	if key == "HTTPPROXY_LISTENING_ADDRESS" {
		return value
	}
	return ":" + value
}

func (r *Reader) readHTTProxyEnabled() (enabled *bool, err error) {
	key, s := r.getEnvWithRetro("HTTPPROXY", "PROXY", "TINYPROXY")
	if s == "" {
		return nil, nil //nolint:nilnil
	}

	enabled = new(bool)
	*enabled, err = binary.Validate(s)
	if err != nil {
		return nil, fmt.Errorf("environment variable %s: %w", key, err)
	}

	return enabled, nil
}

func (r *Reader) readHTTProxyLog() (enabled *bool, err error) {
	key, s := r.getEnvWithRetro("HTTPPROXY_LOG", "PROXY_LOG_LEVEL", "TINYPROXY_LOG")
	if s == "" {
		return nil, nil //nolint:nilnil
	}

	var binaryOptions []binary.Option
	if key != "HTTPROXY_LOG" {
		retroOption := binary.OptionEnabled("on", "info", "connect", "notice")
		binaryOptions = append(binaryOptions, retroOption)
	}

	enabled = new(bool)
	*enabled, err = binary.Validate(s, binaryOptions...)
	if err != nil {
		return nil, fmt.Errorf("environment variable %s: %w", key, err)
	}

	return enabled, nil
}
