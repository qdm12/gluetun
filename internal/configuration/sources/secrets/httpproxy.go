package secrets

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func readHTTPProxy() (settings settings.HTTPProxy, err error) {
	settings.User, err = readSecretFileAsStringPtr(
		"HTTPPROXY_USER_SECRETFILE",
		"/run/secrets/httpproxy_user",
	)
	if err != nil {
		return settings, fmt.Errorf("cannot read HTTP proxy user secret file: %w", err)
	}

	settings.Password, err = readSecretFileAsStringPtr(
		"HTTPPROXY_PASSWORD_SECRETFILE",
		"/run/secrets/httpproxy_password",
	)
	if err != nil {
		return settings, fmt.Errorf("cannot read OpenVPN password secret file: %w", err)
	}

	return settings, nil
}
