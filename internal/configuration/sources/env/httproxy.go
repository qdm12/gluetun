package env

import (
	"fmt"
	"os"
	"strings"

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
	s := os.Getenv("HTTPPROXY_USER")
	if s != "" {
		return &s
	}

	// Retro-compatibility
	s = os.Getenv("TINYPROXY_USER")
	if s != "" {
		r.onRetroActive("TINYPROXY_USER", "HTTPPROXY_USER")
		return &s
	}

	// Retro-compatibility
	s = os.Getenv("PROXY_USER")
	if s != "" {
		r.onRetroActive("PROXY_USER", "HTTPPROXY_USER")
		return &s
	}

	return nil
}

func (r *Reader) readHTTProxyPassword() (user *string) {
	s := os.Getenv("HTTPPROXY_PASSWORD")
	if s != "" {
		return &s
	}

	// Retro-compatibility
	s = os.Getenv("TINYPROXY_PASSWORD")
	if s != "" {
		r.onRetroActive("TINYPROXY_PASSWORD", "HTTPPROXY_PASSWORD")
		return &s
	}

	// Retro-compatibility
	s = os.Getenv("PROXY_PASSWORD")
	if s != "" {
		r.onRetroActive("PROXY_PASSWORD", "HTTPPROXY_PASSWORD")
		return &s
	}

	return nil
}

func (r *Reader) readHTTProxyListeningAddress() (listeningAddress string) {
	// Retro-compatibility
	retroKeys := []string{"PROXY_PORT", "TINYPROXY_PORT", "HTTPPROXY_PORT"}
	for _, retroKey := range retroKeys {
		s := os.Getenv(retroKey)
		if s != "" {
			r.onRetroActive(retroKey, "HTTPPROXY_LISTENING_ADDRESS")
			return ":" + s
		}
	}

	return os.Getenv("HTTPPROXY_LISTENING_ADDRESS")
}

func (r *Reader) readHTTProxyEnabled() (enabled *bool, err error) {
	// Retro-compatibility
	s := strings.ToLower(os.Getenv("PROXY"))
	if s != "" {
		r.onRetroActive("PROXY", "HTTPPROXY")
		enabled = new(bool)
		*enabled, err = binary.Validate(s)
		if err != nil {
			return nil, fmt.Errorf("environment variable PROXY: %w", err)
		}
		return enabled, nil
	}

	// Retro-compatibility
	s = strings.ToLower(os.Getenv("TINYPROXY"))
	if s != "" {
		r.onRetroActive("TINYPROXY", "HTTPPROXY")
		enabled = new(bool)
		*enabled, err = binary.Validate(s)
		if err != nil {
			return nil, fmt.Errorf("environment variable TINYPROXY: %w", err)
		}
		return enabled, nil
	}

	s = strings.ToLower(os.Getenv("HTTPPROXY"))
	if s != "" {
		enabled = new(bool)
		*enabled, err = binary.Validate(s)
		if err != nil {
			return nil, fmt.Errorf("environment variable HTTPPROXY: %w", err)
		}
		return enabled, nil
	}

	return nil, nil //nolint:nilnil
}

func (r *Reader) readHTTProxyLog() (enabled *bool, err error) {
	// Retro-compatibility
	retroOption := binary.OptionEnabled("on", "info", "connect", "notice")
	s := strings.ToLower(os.Getenv("PROXY_LOG_LEVEL"))
	if s != "" {
		r.onRetroActive("PROXY_LOG_LEVEL", "HTTPPROXY_LOG")
		enabled = new(bool)
		*enabled, err = binary.Validate(s, retroOption)
		if err != nil {
			return nil, fmt.Errorf("environment variable PROXY_LOG_LEVEL: %w", err)
		}
		return enabled, nil
	}

	// Retro-compatibility
	s = strings.ToLower(os.Getenv("TINYPROXY_LOG"))
	if s != "" {
		r.onRetroActive("TINYPROXY_LOG", "HTTPPROXY_LOG")
		enabled = new(bool)
		*enabled, err = binary.Validate(s, retroOption)
		if err != nil {
			return nil, fmt.Errorf("environment variable TINYPROXY_LOG: %w", err)
		}
		return enabled, nil
	}

	s = strings.ToLower(os.Getenv("HTTPPROXY_LOG"))
	if s != "" {
		enabled = new(bool)
		*enabled, err = binary.Validate(s)
		if err != nil {
			return nil, fmt.Errorf("environment variable HTTPPROXY_LOG: %w", err)
		}
		return enabled, nil
	}

	return nil, nil //nolint:nilnil
}
