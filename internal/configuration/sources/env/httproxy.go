package env

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/govalid/binary"
)

func (s *Source) readHTTPProxy() (httpProxy settings.HTTPProxy, err error) {
	httpProxy.User = s.readHTTProxyUser()
	httpProxy.Password = s.readHTTProxyPassword()
	httpProxy.ListeningAddress = s.readHTTProxyListeningAddress()

	httpProxy.Enabled, err = s.readHTTProxyEnabled()
	if err != nil {
		return httpProxy, err
	}

	httpProxy.Stealth, err = envToBoolPtr("HTTPPROXY_STEALTH")
	if err != nil {
		return httpProxy, fmt.Errorf("environment variable HTTPPROXY_STEALTH: %w", err)
	}

	httpProxy.Log, err = s.readHTTProxyLog()
	if err != nil {
		return httpProxy, err
	}

	return httpProxy, nil
}

func (s *Source) readHTTProxyUser() (user *string) {
	_, value := s.getEnvWithRetro("HTTPPROXY_USER", "PROXY_USER", "TINYPROXY_USER")
	if value != "" {
		return &value
	}
	return nil
}

func (s *Source) readHTTProxyPassword() (user *string) {
	_, value := s.getEnvWithRetro("HTTPPROXY_PASSWORD", "PROXY_PASSWORD", "TINYPROXY_PASSWORD")
	if value != "" {
		return &value
	}
	return nil
}

func (s *Source) readHTTProxyListeningAddress() (listeningAddress string) {
	key, value := s.getEnvWithRetro("HTTPPROXY_LISTENING_ADDRESS", "PROXY_PORT", "TINYPROXY_PORT", "HTTPPROXY_PORT")
	if key == "HTTPPROXY_LISTENING_ADDRESS" {
		return value
	}
	return ":" + value
}

func (s *Source) readHTTProxyEnabled() (enabled *bool, err error) {
	key, value := s.getEnvWithRetro("HTTPPROXY", "PROXY", "TINYPROXY")
	if value == "" {
		return nil, nil //nolint:nilnil
	}

	enabled = new(bool)
	*enabled, err = binary.Validate(value)
	if err != nil {
		return nil, fmt.Errorf("environment variable %s: %w", key, err)
	}

	return enabled, nil
}

func (s *Source) readHTTProxyLog() (enabled *bool, err error) {
	key, value := s.getEnvWithRetro("HTTPPROXY_LOG", "PROXY_LOG_LEVEL", "TINYPROXY_LOG")
	if value == "" {
		return nil, nil //nolint:nilnil
	}

	var binaryOptions []binary.Option
	if key != "HTTPROXY_LOG" {
		retroOption := binary.OptionEnabled("on", "info", "connect", "notice")
		binaryOptions = append(binaryOptions, retroOption)
	}

	enabled = new(bool)
	*enabled, err = binary.Validate(value, binaryOptions...)
	if err != nil {
		return nil, fmt.Errorf("environment variable %s: %w", key, err)
	}

	return enabled, nil
}
