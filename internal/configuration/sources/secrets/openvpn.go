package secrets

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func readOpenVPN() (
	settings settings.OpenVPN, err error) {
	settings.User, err = readSecretFileAsStringPtr(
		"OPENVPN_USER_SECRETFILE",
		"/run/secrets/openvpn_user",
	)
	if err != nil {
		return settings, fmt.Errorf("cannot read user file: %w", err)
	}

	settings.Password, err = readSecretFileAsStringPtr(
		"OPENVPN_PASSWORD_SECRETFILE",
		"/run/secrets/openvpn_password",
	)
	if err != nil {
		return settings, fmt.Errorf("cannot read password file: %w", err)
	}

	settings.Key, err = readPEMSecretFile(
		"OPENVPN_CLIENTKEY_SECRETFILE",
		"/run/secrets/openvpn_clientkey",
	)
	if err != nil {
		return settings, fmt.Errorf("cannot read client key file: %w", err)
	}

	settings.EncryptedKey, err = readPEMSecretFile(
		"OPENVPN_ENCRYPTED_KEY_SECRETFILE",
		"/run/secrets/openvpn_encrypted_key",
	)
	if err != nil {
		return settings, fmt.Errorf("reading encrypted key file: %w", err)
	}

	settings.KeyPassphrase, err = readSecretFileAsStringPtr(
		"OPENVPN_KEY_PASSPHRASE_SECRETFILE",
		"/run/secrets/openvpn_key_passphrase",
	)
	if err != nil {
		return settings, fmt.Errorf("reading key passphrase file: %w", err)
	}

	settings.Cert, err = readPEMSecretFile(
		"OPENVPN_CLIENTCRT_SECRETFILE",
		"/run/secrets/openvpn_clientcrt",
	)
	if err != nil {
		return settings, fmt.Errorf("cannot read client certificate file: %w", err)
	}

	return settings, nil
}
