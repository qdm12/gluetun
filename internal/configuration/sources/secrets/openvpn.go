package secrets

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func readOpenVPN() (
	settings settings.OpenVPN, err error) {
	settings.User, err = readSecretFileAsString(
		"OPENVPN_USER_SECRETFILE",
		"/run/secrets/openvpn_user",
	)
	if err != nil {
		return settings, fmt.Errorf("cannot read user file: %w", err)
	}

	settings.Password, err = readSecretFileAsString(
		"OPENVPN_PASSWORD_SECRETFILE",
		"/run/secrets/openvpn_password",
	)
	if err != nil {
		return settings, fmt.Errorf("cannot read password file: %w", err)
	}

	settings.ClientKey, err = readSecretFileAsStringPtr(
		"OPENVPN_CLIENTKEY_SECRETFILE",
		"/run/secrets/openvpn_clientkey",
	)
	if err != nil {
		return settings, fmt.Errorf("cannot read client key file: %w", err)
	}

	settings.ClientCrt, err = readSecretFileAsStringPtr(
		"OPENVPN_CLIENTCRT_SECRETFILE",
		"/run/secrets/openvpn_clientcrt",
	)
	if err != nil {
		return settings, fmt.Errorf("cannot read client certificate file: %w", err)
	}

	return settings, nil
}
