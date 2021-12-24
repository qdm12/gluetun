package secrets

import (
	"fmt"
	"os"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/configuration/sources/files"
)

func (r *Reader) readHTTPProxy() (
	settings settings.HTTPProxy, err error) {
	settings.User, err = r.readHTTPProxyUser()
	if err != nil {
		return settings, fmt.Errorf("cannot read HTTP proxy user secret file: %w", err)
	}

	settings.Password, err = r.readHTTPProxyPassword()
	if err != nil {
		return settings, fmt.Errorf("cannot read OpenVPN password secret file: %w", err)
	}

	return settings, nil
}

func (r *Reader) readHTTPProxyUser() (user *string, err error) {
	const envKey = "HTTPPROXY_USER_SECRETFILE"
	path := os.Getenv(envKey)
	if path == "" {
		path = "/run/secrets/httpproxy_user"
	}
	user, err = files.ReadFromFile(path)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *Reader) readHTTPProxyPassword() (password *string, err error) {
	const envKey = "HTTPPROXY_PASSWORD_SECRETFILE"
	path := os.Getenv(envKey)
	if path == "" {
		path = "/run/secrets/httpproxy_password"
	}
	password, err = files.ReadFromFile(path)
	if err != nil {
		return nil, err
	}
	return password, nil
}
